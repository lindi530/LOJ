package upload_service

import (
	"GO1/database/redis/course_redis"
	"GO1/middlewares/response"
	"GO1/models/upload_model"

	"github.com/gin-gonic/gin"
)

func CheckChunkMD5(c *gin.Context, req *upload_model.CheckChunkMD5Req) (resp response.Response) {
	exist, err := course_redis.CheckChunkMD5(c.Request.Context(), req.UploadID, req.ChunkID, req.MD5)
	if err != nil {
		resp.Code = 1
		resp.Message = "分块MD5查询失败"
		return
	}

	resp.Data = &upload_model.CheckChunkMD5Resp{
		Exist: exist,
	}
	return
}
