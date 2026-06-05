package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func GetCompetitionHistory(userID int64) (data []competition_model.CompetitionResult, err error) {
	err = global.DB.
		Select("competition_id, competition_name, competition_rank, rating_before, rating_after, created_at").
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Order("id ASC").
		Find(&data).Error

	return
}
