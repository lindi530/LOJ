package saber_model

import (
	"GO1/pkg/constants"
	"time"
)

func CreateSaberStat(userid int64) SaberStat {
	saberStat := SaberStat{
		UserID:       userid,
		Rating:       constants.SaberInitialRating,
		Level:        constants.SaberLevelGold,
		Wins:         0,
		TotalMatches: 0,
		CreatedAT:    time.Now(),
		UpdatedAT:    time.Now(),
	}
	return saberStat
}
