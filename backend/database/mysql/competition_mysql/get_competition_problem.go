package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
)

func GetCompetitionProblems(competitionID int64, userID int64) (data []competition_model.GetCompetitionProblemsResp, err error) {
	err = global.DB.
		Table("competition_problems cp").
		Select(`
			cp.problem_number,
			cp.problem_id,
			p.title AS problem_title,
			cp.ac_count,
			cp.submit_count,
			EXISTS (
				SELECT 1
				FROM competition_submissions cs
				WHERE cs.competition_id = cp.competition_id
					AND cs.problem_id = cp.problem_id
					AND cs.is_ac = 1
					AND cs.user_id = ?
			) AS has_access
		`, userID).
		Joins("JOIN problems p ON p.id = cp.problem_id").
		Where("cp.competition_id = ?", competitionID).
		Order("cp.problem_number ASC").
		Scan(&data).Error

	return
}
