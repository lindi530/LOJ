package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"

	"gorm.io/gorm"
)

func CreateCompetition(competition *competition_model.Competition, problemIDs []int) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(competition).Error; err != nil {
			return err
		}

		for _, problemID := range problemIDs {
			if err := tx.Create(&competition_model.CompetitionProblem{
				CompetitionID: competition.ID,
				ProblemID:     problemID,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
