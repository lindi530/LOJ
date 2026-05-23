package calendar_service

import (
	"GO1/database/mysql/calendar_mysql"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/calendar_model"
	"GO1/pkg/constants"

	"github.com/gin-gonic/gin"
)

func GetSubmissionCalendar(c *gin.Context, req *calendar_model.CalendarSubmissionReq) (resp response.Response) {
	respData := &calendar_model.CalendarSubmissionResp{}
	startDate := req.StartDate
	endDate := req.EndDate
	userID, exists := c.Get(constants.UserID)
	if !exists {
		resp.Code = 1
		resp.Message = "用户ID错误！"
		return
	}
	global.Logger.Error("UserID: ", userID)
	global.Logger.Error("StartTime: ", startDate)
	global.Logger.Error("EndTime: ", endDate)
	data, err := calendar_mysql.GetSubmissionCalendar(userID.(int64), startDate, endDate)
	if err != nil {
		resp.Code = 1
		resp.Message = "数据查询错误"
		return
	}

	for _, item := range data {
		respData.TotalAC += item.ACCount
		respData.ActiveDays += 1
		respData.Days = append(respData.Days, calendar_model.SubmissionDays{
			Date:    item.ACDate,
			ACCount: item.ACCount,
		})
	}
	resp.Data = respData
	return
}
