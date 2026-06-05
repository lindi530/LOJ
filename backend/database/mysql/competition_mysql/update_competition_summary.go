package competition_mysql

import (
	"GO1/models/competition_model"

	"gorm.io/gorm"
)

func UpdateCompetitionSummary(tx *gorm.DB, results []competition_model.CompetitionResult) error {
	if len(results) == 0 {
		return nil
	}

	userIDs := make([]int64, 0, len(results))
	for _, result := range results {
		userIDs = append(userIDs, result.UserId)
	}

	var existingResults []struct {
		UserID int64 `gorm:"column:user_id"`
	}
	if err := tx.
		Table("competition_results").
		Select("user_id").
		Where("competition_id = ? AND user_id IN ?", results[0].CompetitionId, userIDs).
		Find(&existingResults).Error; err != nil {
		return err
	}

	existingResultUserIDs := make(map[int64]bool, len(existingResults))
	for _, result := range existingResults {
		existingResultUserIDs[result.UserID] = true
	}

	for _, result := range results {
		competitionCountDelta := 1
		if existingResultUserIDs[result.UserId] {
			competitionCountDelta = 0
		}

		if err := tx.Exec(`
			INSERT INTO competition_summary (user_id, rank_score, competition_count)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE
				rank_score = VALUES(rank_score),
				competition_count = competition_count + VALUES(competition_count)
		`, result.UserId, result.RatingAfter, competitionCountDelta).Error; err != nil {
			return err
		}
	}

	return nil
}
