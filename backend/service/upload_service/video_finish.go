package upload_service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"GO1/database/mysql/course_mysql"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/upload_model"

	"github.com/Eyevinn/mp4ff/mp4"
	"github.com/gin-gonic/gin"
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

	chunkDir := filepath.Join(global.Config.Upload.Video.Path, strconv.FormatInt(req.UploadID, 10))
	chunkPaths, err := collectVideoChunkPaths(chunkDir, req.ChunkCount)
	if err != nil {
		resp.Code = 1
		resp.Message = "视频分块不完整"
		return
	}

	finalPath := filepath.Join(chunkDir, filename)
	if hasSamePath(finalPath, chunkPaths) {
		resp.Code = 1
		resp.Message = "视频文件名错误"
		return
	}

	tempPath, fileMD5, sizeBytes, err := mergeVideoChunks(chunkDir, chunkPaths)
	if err != nil {
		resp.Code = 1
		resp.Message = "视频分块合并失败"
		return
	}

	if !strings.EqualFold(fileMD5, req.MD5) {
		_ = os.Remove(tempPath)
		resp.Code = 1
		resp.Message = "视频MD5校验失败"
		return
	}

	if err := replaceFile(tempPath, finalPath); err != nil {
		_ = os.Remove(tempPath)
		resp.Code = 1
		resp.Message = "视频文件保存失败"
		return
	}

	//if err := removeVideoChunks(chunkPaths, finalPath); err != nil {
	//	resp.Code = 1
	//	resp.Message = "视频分块清理失败"
	//	return
	//}

	if err := ensureVideoFile(finalPath, sizeBytes); err != nil {
		resp.Code = 1
		resp.Message = "视频文件保存失败"
		return
	}

	duration := 0
	if parsedDuration, err := readVideoDurationSeconds(finalPath); err == nil {
		duration = parsedDuration
	}

	videoURL := path.Join(global.Config.Upload.Video.Path, strconv.FormatInt(req.UploadID, 10), filename)
	if err := course_mysql.UpdateVideoAssetFinish(req.UploadID, videoURL, sizeBytes, duration); err != nil {
		resp.Code = 1
		resp.Message = "视频资源信息更新失败"
		return
	}

	resp.Data = &upload_model.VideoFinishResp{
		URL: videoURL,
		ID:  req.UploadID,
	}
	return
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

func collectVideoChunkPaths(chunkDir string, chunkCount int) ([]string, error) {
	entries, err := os.ReadDir(chunkDir)
	if err != nil {
		return nil, err
	}

	chunkPaths := make([]string, chunkCount)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		chunkID, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if chunkID < 0 || chunkID >= chunkCount {
			return nil, fmt.Errorf("unexpected chunk id %d", chunkID)
		}

		chunkPaths[chunkID] = filepath.Join(chunkDir, entry.Name())
	}

	for chunkID, chunkPath := range chunkPaths {
		if chunkPath == "" {
			return nil, fmt.Errorf("missing chunk id %d", chunkID)
		}
	}

	return chunkPaths, nil
}

func mergeVideoChunks(chunkDir string, chunkPaths []string) (string, string, int64, error) {
	tempFile, err := os.CreateTemp(chunkDir, ".merge-*")
	if err != nil {
		return "", "", 0, err
	}

	tempPath := tempFile.Name()
	hash := md5.New()
	writer := io.MultiWriter(tempFile, hash)

	var sizeBytes int64
	for _, chunkPath := range chunkPaths {
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			_ = tempFile.Close()
			_ = os.Remove(tempPath)
			return "", "", 0, err
		}

		written, copyErr := io.Copy(writer, chunkFile)
		closeErr := chunkFile.Close()
		if copyErr != nil {
			_ = tempFile.Close()
			_ = os.Remove(tempPath)
			return "", "", 0, copyErr
		}
		if closeErr != nil {
			_ = tempFile.Close()
			_ = os.Remove(tempPath)
			return "", "", 0, closeErr
		}

		sizeBytes += written
	}

	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return "", "", 0, err
	}

	return tempPath, hex.EncodeToString(hash.Sum(nil)), sizeBytes, nil
}

func replaceFile(sourcePath, targetPath string) error {
	if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Rename(sourcePath, targetPath)
}

func removeVideoChunks(chunkPaths []string, keepPath string) error {
	for _, chunkPath := range chunkPaths {
		if hasSamePath(keepPath, []string{chunkPath}) {
			continue
		}
		if err := os.Remove(chunkPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func ensureVideoFile(filePath string, sizeBytes int64) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("video path is directory")
	}
	if info.Size() != sizeBytes {
		return fmt.Errorf("video size mismatch")
	}
	return nil
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

func hasSamePath(targetPath string, paths []string) bool {
	targetPath = filepath.Clean(targetPath)
	for _, item := range paths {
		if filepath.Clean(item) == targetPath {
			return true
		}
	}
	return false
}
