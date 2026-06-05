package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func GetCompetitionUserRatings(userIDs []int64) (map[int64]int, error) {
	ratings := make(map[int64]int, len(userIDs))
	if len(userIDs) == 0 {
		return ratings, nil
	}

	var data []competition_model.CompetitionSummary
	if err := global.DB.
		Table("competition_summary").
		Select("user_id, rank_score").
		Where("user_id IN ?", userIDs).
		Find(&data).Error; err != nil {
		return nil, err
	}

	for _, userID := range userIDs {
		ratings[userID] = 1000
	}
	for _, item := range data {
		ratings[item.UserID] = item.RankScore
	}

	return ratings, nil
}
