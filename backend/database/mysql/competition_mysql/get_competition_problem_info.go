package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
	"GO1/models/problem_model"

	"gorm.io/gorm"
)

func GetCompetitionProblemInfo(competitionID int64, problemNumber string) (data competition_model.GetCompetitionProblemInfoResp, constraints []problem_model.ProblemConstraint, examples []problem_model.Example, err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		var competitionProblem competition_model.CompetitionProblem
		if err := tx.
			Table("competition_problems").
			Select("problem_id").
			Where("competition_id = ? AND problem_number = ?", competitionID, problemNumber).
			Take(&competitionProblem).Error; err != nil {
			return err
		}

		if err := tx.
			Table("problems").
			Select("title, description, input_desc, output_desc").
			Where("id = ?", competitionProblem.ProblemID).
			Take(&data).Error; err != nil {
			return err
		}

		constraints, err = getCompetitionProblemConstraints(tx, uint(competitionProblem.ProblemID))
		if err != nil {
			return err
		}

		examples, err = getCompetitionProblemExamples(tx, uint(competitionProblem.ProblemID))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return data, nil, nil, err
	}

	return data, constraints, examples, nil
}

func getCompetitionProblemConstraints(tx *gorm.DB, problemID uint) (data []problem_model.ProblemConstraint, err error) {
	err = tx.
		Select("time_limit, memory_limit, language").
		Where("problem_id = ?", problemID).
		Find(&data).Error

	return
}

func getCompetitionProblemExamples(tx *gorm.DB, problemID uint) (data []problem_model.Example, err error) {
	err = tx.
		Model(&problem_model.ProblemExample{}).
		Select("input, output, is_public").
		Where("problem_id = ? AND is_public = ?", problemID, true).
		Find(&data).Error

	return
}
