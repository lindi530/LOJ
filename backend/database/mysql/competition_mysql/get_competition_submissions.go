package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"

	"gorm.io/gorm"
)

func GetCompetitionSubmissions(competitionID int64, offset, limit int, problemNumber string) (data competition_model.GetCompetitionSubmissionsResp, err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := getCompetitionSubmissionsBaseQuery(tx, competitionID, problemNumber).
			Count(&data.Count).Error; err != nil {
			return err
		}

		return getCompetitionSubmissionsBaseQuery(tx, competitionID, problemNumber).
			Select(`
				ps.id AS submission_id,
				cs.user_id AS submission_user_id,
				COALESCE(cpl.user_name, '') AS submission_user_name,
				COALESCE(cp.problem_number, '') AS problem_number,
				COALESCE(p.title, '') AS problem_title,
				ps.state,
				ps.language,
				ps.exec_time AS exec_limit,
				ps.memory_usage AS memory_limit,
				ps.created_at AS submission_time
			`).
			Joins("LEFT JOIN competition_players AS cpl ON cpl.competition_id = cs.competition_id AND cpl.user_id = cs.user_id").
			Joins("LEFT JOIN problems AS p ON p.id = cs.problem_id").
			Order("ps.created_at DESC, ps.id DESC").
			Offset(offset).
			Limit(limit).
			Scan(&data.List).Error
	})

	return
}

func getCompetitionSubmissionsBaseQuery(tx *gorm.DB, competitionID int64, problemNumber string) *gorm.DB {
	query := tx.
		Table("competition_submissions AS cs").
		Joins("JOIN problem_submissions AS ps ON ps.id = cs.submission_id").
		Joins("LEFT JOIN competition_problems AS cp ON cp.competition_id = cs.competition_id AND cp.problem_id = cs.problem_id").
		Where("cs.competition_id = ?", competitionID)

	if problemNumber != "" {
		query = query.Where("cp.problem_number = ?", problemNumber)
	}

	return query
}
