package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
	"GO1/models/problem_model"
	"GO1/pkg/constants"

	"gorm.io/gorm"
)

func GetCompetitionProblemForSubmit(competitionID int64, problemNumber string) (competition competition_model.Competition, competitionProblem competition_model.CompetitionProblem, examples []problem_model.Example, err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Select("id, start_time, end_time").
			Where("id = ?", competitionID).
			Take(&competition).Error; err != nil {
			return err
		}

		if err := tx.
			Where("competition_id = ? AND problem_number = ?", competitionID, problemNumber).
			Take(&competitionProblem).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&problem_model.ProblemExample{}).
			Select("input, output").
			Where("problem_id = ?", competitionProblem.ProblemID).
			Find(&examples).Error; err != nil {
			return err
		}

		return nil
	})

	return
}

func HasAcceptedCompetitionProblem(competitionID int64, problemID int, userID int64) (bool, error) {
	var count int64
	err := global.DB.
		Table("competition_submissions AS cs").
		Joins("JOIN problem_submissions AS ps ON ps.id = cs.submission_id").
		Where("cs.competition_id = ? AND cs.problem_id = ? AND cs.user_id = ? AND ps.state = ?", competitionID, problemID, userID, constants.JudgeStatusAccepted).
		Count(&count).Error

	return count > 0, err
}
