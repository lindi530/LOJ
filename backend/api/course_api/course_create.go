package course_api

import (
	"GO1/middlewares/response"
	"GO1/models/course_model"
	"GO1/pkg/jwt"
	"GO1/service/course_service"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) CourseCreate(c *gin.Context) {
	var req course_model.CourseCreateReq

	jwt.SaveUserInfoFromToken(c)

	if err := c.ShouldBind(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	cover, err := c.FormFile("cover")
	if err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	resp := course_service.CourseCreate(c, &req, cover)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
