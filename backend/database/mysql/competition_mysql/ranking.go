package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func GetRanking() (data []competition_model.GetRankingResp, err error) {
	err = global.DB.
		Table("competition_summary AS cs").
		Select(`
			cs.user_id,
			u.user_name,
			u.avatar AS user_avatar,
			cs.rank_score AS user_rating,
			cs.competition_count AS user_competition_count
		`).
		Joins("JOIN users AS u ON u.user_id = cs.user_id").
		Order("cs.rank_score DESC").
		Order("cs.competition_count DESC").
		Order("cs.user_id ASC").
		Scan(&data).Error

	return
}
