package course_api

import (
	"GO1/middlewares/response"
	"GO1/models/course_model"
	service_context "GO1/service/context"
	"GO1/service/course_service"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) GetChapterInfo(c *gin.Context) {
	var req course_model.GetChapterInfoReq

	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	userID, ok := service_context.GetContext(c, service_context.CtxUserIdKey).(int64)
	if !ok || userID <= 0 {
		response.FailWithCode(response.NeedLogin, c)
		return
	}

	resp := course_service.ChapterInfo(userID, &req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
