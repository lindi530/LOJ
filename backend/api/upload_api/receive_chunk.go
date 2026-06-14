package upload_api

import (
	"strconv"

	"GO1/middlewares/response"
	"GO1/models/upload_model"
	"GO1/service/upload_service"

	"github.com/gin-gonic/gin"
)

func (UploadAPI) ReceiveChunk(c *gin.Context) {
	var uriReq struct {
		UploadID int64 `uri:"upload_id" binding:"required,min=1"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	chunkID, err := strconv.Atoi(c.PostForm("chunk_id"))
	if err != nil || chunkID < 0 {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	chunk, err := c.FormFile("chunk")
	if err != nil || chunk == nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	req := upload_model.ReceiveChunkReq{
		UploadID: uriReq.UploadID,
		ChunkID:  chunkID,
		Chunk:    chunk,
	}

	resp := upload_service.ReceiveChunk(c, &req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
