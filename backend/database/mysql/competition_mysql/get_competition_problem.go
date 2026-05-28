package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func GetCompetitionTime(competitionID int64) (data competition_model.Competition, err error) {
	err = global.DB.
		Select("id, start_time, end_time").
		Where("id = ?", competitionID).
		First(&data).Error

	return
}

func GetCompetitionProblems(competitionID int64) (data []competition_model.GetCompetitionProblemsResp, err error) {
	err = global.DB.
		Table("competition_problems cp").
		Select("cp.problem_number, cp.problem_id, p.title AS problem_title, cp.ac_count, cp.submit_count").
		Joins("JOIN problems p ON p.id = cp.problem_id").
		Where("cp.competition_id = ?", competitionID).
		Order("cp.problem_number ASC").
		Scan(&data).Error

	return
}
