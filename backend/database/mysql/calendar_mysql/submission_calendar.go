package calendar_mysql

import (
	"GO1/global"
	"GO1/models/calendar_model"
	"time"
)

func GetSubmissionCalendar(userID int64, startTime, endTime time.Time) (data []calendar_model.CalendarSubmission, err error) {

	err = global.DB.
		Select("ac_count, ac_date").
		Where("user_id = ? AND ac_date >= ? AND ac_date <= ?", userID, startTime, endTime).
		Order("ac_date ASC").
		Find(&data).Error

	global.Logger.Error("数据库错误：", err)
	return
}
