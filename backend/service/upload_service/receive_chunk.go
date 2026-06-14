package upload_service

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"GO1/database/redis/course_redis"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/upload_model"

	"github.com/gin-gonic/gin"
)

func ReceiveChunk(c *gin.Context, req *upload_model.ReceiveChunkReq) (resp response.Response) {
	if req.Chunk == nil {
		resp.Code = 1
		resp.Message = "视频分块不能为空"
		return
	}

	chunkMD5, err := calculateChunkMD5(req.Chunk)
	if err != nil {
		resp.Code = 1
		resp.Message = "视频分块MD5校验失败"
		return
	}

	if err := saveChunkFile(c, req); err != nil {
		resp.Code = 1
		resp.Message = "视频分块保存失败"
		return
	}

	if err := course_redis.SaveChunkMD5(c.Request.Context(), req.UploadID, req.ChunkID, chunkMD5); err != nil {
		resp.Code = 1
		resp.Message = "视频分块MD5保存失败"
		return
	}

	resp.Data = &upload_model.ReceiveChunkResp{}
	return
}

func calculateChunkMD5(chunk *multipart.FileHeader) (string, error) {
	file, err := chunk.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func saveChunkFile(c *gin.Context, req *upload_model.ReceiveChunkReq) error {
	chunkDir := filepath.Join(global.Config.Upload.Video.Path, strconv.FormatInt(req.UploadID, 10))
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return err
	}

	chunkPath := filepath.Join(chunkDir, strconv.Itoa(req.ChunkID))
	return c.SaveUploadedFile(req.Chunk, chunkPath)
}
