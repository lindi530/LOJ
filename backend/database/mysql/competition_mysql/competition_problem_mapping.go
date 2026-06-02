package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func GetCompetitionProblemMappings(competitionID int64) (data []competition_model.CompetitionProblem, err error) {
	err = global.DB.
		Table("competition_problems").
		Select("problem_number, problem_id").
		Where("competition_id = ?", competitionID).
		Order("problem_number ASC").
		Find(&data).Error

	return
}
