package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func GetCompetitionTime(competitionID int64) (data competition_model.Competition, err error) {
	err = global.DB.
		Select("id, start_time, end_time").
		Where("id = ?", competitionID).
		First(&data).Error

	return
}
