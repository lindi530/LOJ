package routers

import (
	"GO1/api"
	"GO1/middlewares/jwt"

	"github.com/gin-gonic/gin"
)

func UploadRouter(router *gin.RouterGroup) {
	upload := router.Group("/upload")

	upload.Use(jwt.JWTAuthMiddleware())
	{
		upload.POST("problem", api.ApiGroups.UploadAPI.UploadProblem)

	}

	video := upload.Group("/video")
	video.Use(jwt.JWTAuthMiddleware())
	{
		video.POST("", api.ApiGroups.UploadAPI.CreateVideoUploadTask)         // 创建视频上传任务，返回upload_id 和 chunk_size
		video.GET("/:upload_id", api.ApiGroups.UploadAPI.CheckChunkMD5)       // 根据md5值判断该分块是否存在
		video.POST("/:upload_id", api.ApiGroups.UploadAPI.ReceiveChunk)       // 上传每一个chunk
		video.POST("/:upload_id/finish", api.ApiGroups.UploadAPI.VideoFinish) // 最终合并
	}
}
