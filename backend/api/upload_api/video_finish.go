package upload_api

import (
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/upload_model"
	"GO1/service/upload_service"

	"github.com/gin-gonic/gin"
)

func (UploadAPI) VideoFinish(c *gin.Context) {
	var uriReq struct {
		UploadID int64 `uri:"upload_id" binding:"required,min=1"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	var req upload_model.VideoFinishReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	if req.MD5 == "" || req.FileName == "" || req.ChunkCount <= 0 {
		response.FailWithCode(response.BadRequest, c)
		return
	}
	if req.UploadID != 0 && req.UploadID != uriReq.UploadID {
		response.FailWithCode(response.BadRequest, c)
		return
	}
	req.UploadID = uriReq.UploadID

	resp := upload_service.VideoFinish(c, &req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}
	global.Logger.Warn("resp.Data: ", resp.Data)
	response.OkWithData(resp.Data, c)
}
