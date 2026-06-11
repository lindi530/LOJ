package course_api

import (
	"GO1/middlewares/response"
	"GO1/models/course_model"
	"GO1/service/course_service"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) CourseList(c *gin.Context) {
	var req course_model.CourseListReq

	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	resp := course_service.CourseList(&req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
