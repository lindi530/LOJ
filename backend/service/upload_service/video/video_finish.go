package video

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"strconv"
	"strings"

	"GO1/database/mysql/course_mysql"
	"GO1/database/redis/course_redis"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/upload_model"

	"github.com/Eyevinn/mp4ff/mp4"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	miniogo "github.com/minio/minio-go/v7"
)

func VideoFinish(c *gin.Context, req *upload_model.VideoFinishReq) (resp response.Response) {
	if _, err := course_mysql.GetVideoAssetByID(req.UploadID); err != nil {
		resp.Code = 1
		resp.Message = "视频上传任务不存在"
		return
	}

	filename := safeVideoFilename(req.FileName)
	if filename == "" {
		resp.Code = 1
		resp.Message = "视频文件名错误"
		return
	}
	if req.ChunkCount <= 0 || req.ChunkCount > maxMinioPartCount {
		resp.Code = 1
		resp.Message = "视频分块数量错误"
		return
	}

	uploadMeta, err := course_redis.GetVideoMultipartUpload(c.Request.Context(), req.UploadID)
	if errors.Is(err, redis.Nil) {
		resp.Code = 1
		resp.Message = "视频上传任务不存在或已过期"
		return
	}
	if err != nil {
		resp.Code = 1
		resp.Message = "视频上传任务查询失败"
		return
	}

	parts, sizeBytes, err := collectVideoMultipartParts(c.Request.Context(), req.UploadID, req.ChunkCount)
	if err != nil {
		resp.Code = 1
		resp.Message = "视频分块不完整"
		return
	}

	finalObject := rawVideoObject(req.UploadID, filename)
	if finalObject == uploadMeta.TempObject {
		resp.Code = 1
		resp.Message = "视频文件名错误"
		return
	}

	finalContentType := videoContentType(finalObject)

	if _, err := global.MinioCore.CompleteMultipartUpload(
		c.Request.Context(),
		uploadMeta.Bucket,
		uploadMeta.TempObject,
		uploadMeta.MinioUploadID,
		parts,
		miniogo.PutObjectOptions{ContentType: finalContentType},
	); err != nil {
		resp.Code = 1
		resp.Message = "视频分块合并失败"
		return
	}

	tempPath, fileMD5, downloadedSize, err := downloadMinioObjectWithMD5(c.Request.Context(), uploadMeta.Bucket, uploadMeta.TempObject)
	if err != nil {
		resp.Code = 1
		resp.Message = "视频文件校验失败"
		return
	}
	defer os.Remove(tempPath)

	if downloadedSize != sizeBytes {
		_ = global.MinioCore.RemoveObject(context.Background(), uploadMeta.Bucket, uploadMeta.TempObject, miniogo.RemoveObjectOptions{})
		_ = course_redis.DeleteVideoMultipartUpload(context.Background(), req.UploadID)
		resp.Code = 1
		resp.Message = "视频大小校验失败"
		return
	}

	if !strings.EqualFold(fileMD5, req.MD5) {
		_ = global.MinioCore.RemoveObject(context.Background(), uploadMeta.Bucket, uploadMeta.TempObject, miniogo.RemoveObjectOptions{})
		_ = course_redis.DeleteVideoMultipartUpload(context.Background(), req.UploadID)
		resp.Code = 1
		resp.Message = "视频MD5校验失败"
		return
	}

	duration := 0
	if parsedDuration, err := readVideoDurationSeconds(tempPath); err == nil {
		duration = parsedDuration
	}

	if _, err := global.MinioCore.ComposeObject(
		c.Request.Context(),
		miniogo.CopyDestOptions{
			Bucket:          uploadMeta.Bucket,
			Object:          finalObject,
			ContentType:     finalContentType,
			ReplaceMetadata: true,
		},
		miniogo.CopySrcOptions{
			Bucket: uploadMeta.Bucket,
			Object: uploadMeta.TempObject,
		},
	); err != nil {
		resp.Code = 1
		resp.Message = "视频文件保存失败"
		return
	}

	originPath := path.Join(uploadMeta.Bucket, finalObject)
	if err := course_mysql.UpdateVideoAssetFinish(req.UploadID, originPath, sizeBytes, duration); err != nil {
		resp.Code = 1
		resp.Message = "视频资源信息更新失败"
		return
	}

	if err := sendVideoTranscodeTask(videoTranscodeTask{
		VideoAssetID: req.UploadID,
		Bucket:       uploadMeta.Bucket,
		OriginObject: finalObject,
		OriginPath:   originPath,
	}); err != nil {
		if updateErr := course_mysql.UpdateVideoAssetTranscodeFailed(req.UploadID); updateErr != nil {
			global.Logger.Errorw(
				"update video asset transcode failed status failed",
				"video_asset_id", req.UploadID,
				"err", updateErr,
			)
		}
		resp.Code = 1
		resp.Message = "视频转码任务发送失败"
		return
	}

	_ = global.MinioCore.RemoveObject(context.Background(), uploadMeta.Bucket, uploadMeta.TempObject, miniogo.RemoveObjectOptions{})
	_ = course_redis.DeleteVideoMultipartUpload(context.Background(), req.UploadID)

	resp.Data = &upload_model.VideoFinishResp{
		OriginPath: originPath,
		ID:         req.UploadID,
	}
	return
}

func collectVideoMultipartParts(ctx context.Context, uploadID int64, chunkCount int) ([]miniogo.CompletePart, int64, error) {
	uploadMeta, err := course_redis.GetVideoMultipartUpload(ctx, uploadID)
	if err != nil {
		return nil, 0, err
	}
	uploadedParts, err := listUploadedVideoParts(ctx, uploadMeta)
	if err != nil {
		return nil, 0, err
	}

	parts := make([]miniogo.CompletePart, 0, chunkCount)
	var sizeBytes int64

	for chunkID := 0; chunkID < chunkCount; chunkID++ {
		partNumber := chunkID + 1
		part, exist := uploadedParts[partNumber]
		if !exist {
			return nil, 0, fmt.Errorf("missing chunk %d", chunkID)
		}
		if part.PartNumber != partNumber || part.ETag == "" {
			return nil, 0, fmt.Errorf("invalid chunk %d", chunkID)
		}

		parts = append(parts, objectPartToCompletePart(part))
		sizeBytes += part.Size
	}

	return parts, sizeBytes, nil
}

func downloadMinioObjectWithMD5(ctx context.Context, bucket string, object string) (string, string, int64, error) {
	reader, _, _, err := global.MinioCore.GetObject(ctx, bucket, object, miniogo.GetObjectOptions{})
	if err != nil {
		return "", "", 0, err
	}
	defer reader.Close()

	tempFile, err := os.CreateTemp("", "loj-video-*")
	if err != nil {
		return "", "", 0, err
	}

	tempPath := tempFile.Name()
	hash := md5.New()
	written, copyErr := io.Copy(io.MultiWriter(tempFile, hash), reader)
	closeErr := tempFile.Close()
	if copyErr != nil {
		_ = os.Remove(tempPath)
		return "", "", 0, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return "", "", 0, closeErr
	}

	return tempPath, hex.EncodeToString(hash.Sum(nil)), written, nil
}

func safeVideoFilename(filename string) string {
	filename = strings.TrimSpace(filename)
	filename = strings.ReplaceAll(filename, "\\", "/")
	filename = path.Base(filename)
	if filename == "." || filename == ".." || filename == "/" {
		return ""
	}
	return filename
}

func rawVideoObject(uploadID int64, filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	if ext == "" {
		ext = ".mp4"
	}
	return path.Join(videoSourceObjectPrefix(), strconv.FormatInt(uploadID, 10)+ext)
}

func videoContentType(filename string) string {
	contentType := mime.TypeByExtension(strings.ToLower(path.Ext(filename)))
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}

func readVideoDurationSeconds(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	parsedFile, err := mp4.DecodeFile(file, mp4.WithDecodeMode(mp4.DecModeLazyMdat))
	if err != nil {
		return 0, err
	}
	if parsedFile.Moov == nil || parsedFile.Moov.Mvhd == nil {
		return 0, fmt.Errorf("missing mp4 movie header")
	}

	mvhd := parsedFile.Moov.Mvhd
	return durationSeconds(mvhd.Duration, uint64(mvhd.Timescale))
}

func durationSeconds(duration, timescale uint64) (int, error) {
	if timescale == 0 {
		return 0, fmt.Errorf("invalid mp4 timescale")
	}

	seconds := duration / timescale
	if duration%timescale != 0 {
		seconds++
	}
	if seconds > uint64(^uint(0)>>1) {
		return 0, fmt.Errorf("duration overflow")
	}

	return int(seconds), nil
}
