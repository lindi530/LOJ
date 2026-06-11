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
		upload.POST("")                   // 创建任务，返回upload_id 和 chunk_size
		upload.GET("/:upload_id")         // 返回需要上传的chunk_id
		upload.POST("/:upload_id")        // 上传每一个chunk
		upload.POST("/:upload_id/finesh") // 最终合并
	}
}
