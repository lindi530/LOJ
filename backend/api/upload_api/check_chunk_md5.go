package upload_api

import (
	"GO1/middlewares/response"
	"GO1/models/upload_model"
	"GO1/service/upload_service"

	"github.com/gin-gonic/gin"
)

func (UploadAPI) CheckChunkMD5(c *gin.Context) {
	var uriReq struct {
		UploadID int64 `uri:"upload_id" binding:"required,min=1"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	var queryReq struct {
		ChunkID *int   `form:"chunk_id" binding:"required"`
		MD5     string `form:"md5" binding:"required"`
	}

	if err := c.ShouldBindQuery(&queryReq); err != nil || queryReq.ChunkID == nil || *queryReq.ChunkID < 0 {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	req := upload_model.CheckChunkMD5Req{
		UploadID: uriReq.UploadID,
		ChunkID:  *queryReq.ChunkID,
		MD5:      queryReq.MD5,
	}

	resp := upload_service.CheckChunkMD5(c, &req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
