package calendar_mysql

import (
	"GO1/global"
	"GO1/models/calendar_model"
	"GO1/models/mq_models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveACCount(data *mq_models.ACProblem) error {
	record := calendar_model.CalendarSubmission{
		UserID:  data.UserID,
		ACDate:  data.Date,
		ACCount: data.AcCount,
	}

	return global.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "ac_date"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"ac_count": gorm.Expr("ac_count + ?", data.AcCount),
		}),
	}).Create(&record).Error
}
