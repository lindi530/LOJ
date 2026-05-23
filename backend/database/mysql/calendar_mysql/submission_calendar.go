package calendar_mysql

import (
	"GO1/global"
	"GO1/models/calendar_model"
	"time"
)

func GetSubmissionCalendar(userID int64, startDate, endDate time.Time) (data []calendar_model.CalendarSubmission, err error) {

	err = global.DB.
		Select("ac_count, ac_date").
		Where("user_id = ? AND ac_date >= ? AND ac_date <= ?", userID, startDate, endDate).
		Order("ac_date ASC").
		Find(&data).Error

	return
}
