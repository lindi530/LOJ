package upload_api

import (
	"GO1/middlewares/response"
	"GO1/service/upload_service/video"

	"github.com/gin-gonic/gin"
)

func (UploadAPI) CreateVideoUploadTask(c *gin.Context) {

	resp := video.CreateVideoUploadTask(c.Request.Context())
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
