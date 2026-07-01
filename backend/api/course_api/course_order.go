package course_api

import (
	"GO1/middlewares/response"
	"GO1/models/course_model"
	service_context "GO1/service/context"
	"GO1/service/course_service"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) CreateCourseOrder(c *gin.Context) {
	var req course_model.CreateCourseOrderReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	userID, ok := service_context.GetContext(c, service_context.CtxUserIdKey).(int64)
	if !ok || userID <= 0 {
		response.FailWithCode(response.NeedLogin, c)
		return
	}

	resp := course_service.CreateCourseOrder(userID, &req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.Ok(resp.Data, resp.Message, c)
}
