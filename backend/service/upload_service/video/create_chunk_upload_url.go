package video

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"GO1/database/redis/course_redis"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/upload_model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	maxMinioPartCount       = 10000
	chunkUploadURLExpiresIn = 15 * time.Minute
)

func CreateChunkUploadURL(c *gin.Context, req *upload_model.CreateChunkUploadURLReq) (resp response.Response) {
	if req.ChunkID < 0 || req.ChunkID >= maxMinioPartCount {
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

	partNumber := req.ChunkID + 1
	part, exist, err := getUploadedVideoPart(c.Request.Context(), uploadMeta, partNumber)
	if err != nil {
		resp.Code = 1
		resp.Message = "获取信息失败"
		return
	}
	if exist && objectPartMatchesMD5(part, req.MD5) {
		resp.Data = &upload_model.CreateChunkUploadURLResp{
			Exist:      true,
			PartNumber: partNumber,
		}
		return
	}

	reqParams := url.Values{}
	reqParams.Set("partNumber", strconv.Itoa(partNumber))
	reqParams.Set("uploadId", uploadMeta.MinioUploadID)

	headers := http.Header{}
	respHeaders := map[string]string{}
	if req.MD5 != "" {
		contentMD5, err := md5HexToBase64(req.MD5)
		if err != nil {
			resp.Code = 1
			resp.Message = "视频分块MD5格式错误"
			return
		}
		headers.Set("Content-MD5", contentMD5)
		respHeaders["Content-MD5"] = contentMD5
	}

	uploadURL, err := global.MinioPresignCore.PresignHeader(
		c.Request.Context(),
		http.MethodPut,
		uploadMeta.Bucket,
		uploadMeta.TempObject,
		chunkUploadURLExpiresIn,
		reqParams,
		headers,
	)
	if err != nil {
		global.Logger.Errorw(
			"create chunk upload url failed",
			"upload_id", req.UploadID,
			"chunk_id", req.ChunkID,
			"bucket", uploadMeta.Bucket,
			"object", uploadMeta.TempObject,
			"err", err,
		)
		resp.Code = 1
		resp.Message = "视频分块上传地址创建失败"
		return
	}

	resp.Data = &upload_model.CreateChunkUploadURLResp{
		Exist:      false,
		Method:     http.MethodPut,
		UploadURL:  uploadURL.String(),
		Headers:    respHeaders,
		PartNumber: partNumber,
	}

	return
}

func md5HexToBase64(md5Hex string) (string, error) {
	sum, err := hex.DecodeString(md5Hex)
	if err != nil {
		return "", err
	}
	if len(sum) != 16 {
		return "", errors.New("invalid md5 length")
	}
	return base64.StdEncoding.EncodeToString(sum), nil
}
