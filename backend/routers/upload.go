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
		video.POST("", api.ApiGroups.UploadAPI.CreateVideoUploadTask)
		video.POST("/:upload_id/chunk_url", api.ApiGroups.UploadAPI.CreateChunkUploadURL)
		video.POST("/:upload_id/finish", api.ApiGroups.UploadAPI.VideoFinish)
	}
}
