package course_api

import (
	"GO1/middlewares/response"
	"GO1/models/course_model"
	"GO1/service/course_service"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) GetChapterList(c *gin.Context) {
	var req course_model.GetChapterListReq

	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	resp := course_service.ChapterList(&req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
