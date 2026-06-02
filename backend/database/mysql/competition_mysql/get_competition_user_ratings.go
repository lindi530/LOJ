package competition_mysql

import "GO1/global"

type competitionUserRating struct {
	UserID int64 `gorm:"column:user_id"`
	Rating int   `gorm:"column:rating"`
}

func GetCompetitionUserRatings(userIDs []int64) (map[int64]int, error) {
	ratings := make(map[int64]int, len(userIDs))
	if len(userIDs) == 0 {
		return ratings, nil
	}

	var data []competitionUserRating
	if err := global.DB.
		Table("saber_stats").
		Select("user_id, rating").
		Where("user_id IN ?", userIDs).
		Find(&data).Error; err != nil {
		return nil, err
	}

	for _, userID := range userIDs {
		ratings[userID] = 1500
	}
	for _, item := range data {
		ratings[item.UserID] = item.Rating
	}

	return ratings, nil
}
