package upload_api

import (
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/upload_model"
	"GO1/service/upload_service/video"

	"github.com/gin-gonic/gin"
)

func (UploadAPI) CreateChunkUploadURL(c *gin.Context) {
	var uriReq struct {
		UploadID int64 `uri:"upload_id" binding:"required,min=1"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	var bodyReq struct {
		ChunkID *int   `json:"chunk_id" binding:"required"`
		MD5     string `json:"md5"`
	}
	if err := c.ShouldBindJSON(&bodyReq); err != nil || bodyReq.ChunkID == nil || *bodyReq.ChunkID < 0 {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	req := upload_model.CreateChunkUploadURLReq{
		UploadID: uriReq.UploadID,
		ChunkID:  *bodyReq.ChunkID,
		MD5:      bodyReq.MD5,
	}

	resp := video.CreateChunkUploadURL(c, &req)

	global.Logger.Warnw(
		"chunk_upload_url",
		"code", resp.Code,
		"err_code", resp.ErrCode,
		"data", resp.Data,
		"message", resp.Message,
	)

	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}
	response.OkWithData(resp.Data, c)
}
