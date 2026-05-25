package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func GetCompetitions() (data []*competition_model.Competition, err error) {
	err = global.DB.
		Select("id, name, status, visibility, player_count, start_time, end_time, creator_id, creator_name, created_at, updated_at").
		Find(&data).Error

	return
}
