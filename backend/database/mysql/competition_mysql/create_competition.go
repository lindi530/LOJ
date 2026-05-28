package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"

	"gorm.io/gorm"
)

func CreateCompetition(competition *competition_model.Competition, problems []competition_model.CreateCompetitionProblem) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(competition).Error; err != nil {
			return err
		}

		for _, problem := range problems {
			if err := tx.Create(&competition_model.CompetitionProblem{
				CompetitionID: competition.ID,
				ProblemID:     problem.ProblemID,
				ProblemNumber: problem.ProblemNumber,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
