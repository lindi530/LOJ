package calendar_api

import (
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/calendar_model"
	"GO1/pkg/jwt"
	"GO1/service/calendar_service"

	"github.com/gin-gonic/gin"
)

func (CalendarAPI) GetSubmissionCalendar(c *gin.Context) {

	var scr calendar_model.CalendarSubmissionReq
	if err := c.ShouldBind(&scr); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}
	jwt.SaveUserIDFromToken(c)

	resp := calendar_service.GetSubmissionCalendar(c, &scr)

	global.Logger.Error("resp: ", resp)

	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)

	//global.Logger.Error("日历测试")
	//response.Ok(nil, "", c)
}
