package saber_stats_mysql

import (
	"GO1/global"
	"GO1/models/saber_model"
	"GO1/pkg/constants"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateSaberPKResult(winnerID, loserID int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		winnerStat, err := getSaberStatForUpdate(tx, winnerID)
		if err != nil {
			return err
		}

		loserStat, err := getSaberStatForUpdate(tx, loserID)
		if err != nil {
			return err
		}

		applySaberPKResult(&winnerStat, true)
		applySaberPKResult(&loserStat, false)

		if err := updateSaberStat(tx, &winnerStat); err != nil {
			return err
		}
		return updateSaberStat(tx, &loserStat)
	})
}

func getSaberStatForUpdate(tx *gorm.DB, userID int64) (saber_model.SaberStat, error) {
	now := time.Now()
	stat := saber_model.SaberStat{}
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(saber_model.SaberStat{UserID: userID}).
		Attrs(saber_model.SaberStat{
			Rating:       constants.SaberInitialRating,
			Level:        constants.SaberLevelGold,
			Wins:         0,
			TotalMatches: 0,
			CreatedAT:    now,
			UpdatedAT:    now,
		}).
		FirstOrCreate(&stat).
		Error
	return stat, err
}

func applySaberPKResult(stat *saber_model.SaberStat, win bool) {
	if win {
		stat.Rating += constants.SaberPKRatingDelta
		stat.Wins++
	} else {
		stat.Rating -= constants.SaberPKRatingDelta
		if stat.Rating < 0 {
			stat.Rating = 0
		}
	}

	stat.TotalMatches++
	stat.Level = constants.SaberLevelByRating(stat.Rating)
	stat.UpdatedAT = time.Now()
}

func updateSaberStat(tx *gorm.DB, stat *saber_model.SaberStat) error {
	return tx.
		Model(&saber_model.SaberStat{}).
		Where("user_id = ?", stat.UserID).
		Updates(map[string]interface{}{
			"rating":        stat.Rating,
			"level":         stat.Level,
			"wins":          stat.Wins,
			"total_matches": stat.TotalMatches,
			"updated_at":    stat.UpdatedAT,
		}).
		Error
}
