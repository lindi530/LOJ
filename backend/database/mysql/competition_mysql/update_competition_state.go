package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
	"time"
)

func GetNeedStartCompetitions(now time.Time) (data []competition_model.Competition, err error) {
	err = global.DB.
		Select("id, name, status, start_time, end_time").
		Where("status = ? AND start_time <= ? AND end_time > ?", 1, now, now).
		Find(&data).Error

	return
}

func GetNeedEndCompetitions(now time.Time) (data []competition_model.Competition, err error) {
	err = global.DB.
		Select("id, name, status, start_time, end_time").
		Where("status <> ? AND end_time <= ?", 3, now).
		Find(&data).Error

	return
}

func UpdateCompetitionStatus(competitionID int64, oldStatus []int8, newStatus int8) (bool, error) {
	tx := global.DB.
		Model(&competition_model.Competition{}).
		Where("id = ? AND status IN ?", competitionID, oldStatus).
		Update("status", newStatus)

	return tx.RowsAffected > 0, tx.Error
}
