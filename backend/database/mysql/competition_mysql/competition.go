package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
	"time"
)

func GetCompetitions(hasEnded bool, page, pageSize int) (data []*competition_model.Competition, err error) {
	query := global.DB.Select("*")
	now := time.Now()

	if hasEnded {
		query = query.
			Where("end_time <= ?", now).
			Offset((page - 1) * pageSize).
			Limit(pageSize)
	} else {
		query = query.Where("end_time > ?", now)
	}

	err = query.
		Order("start_time ASC").
		Order("end_time ASC").
		Order("id ASC").
		Find(&data).Error

	return
}
