package course_api

import (
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/course_model"
	"GO1/service/course_service"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) GetChapterVideo(c *gin.Context) {
	var req course_model.GetChapterVideoReq

	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	resp := course_service.GetChapterVideo(&req)

	global.Logger.Info("respMessage: ", resp.Message)
	global.Logger.Info("respData: ", resp.Data)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
