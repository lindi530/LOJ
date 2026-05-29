package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"

	"gorm.io/gorm"
)

func GetCompetitionSubmissions(competitionID int64, offset, limit int) (data competition_model.GetCompetitionSubmissionsResp, err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := getCompetitionSubmissionsBaseQuery(tx, competitionID).
			Count(&data.Count).Error; err != nil {
			return err
		}

		return getCompetitionSubmissionsBaseQuery(tx, competitionID).
			Select(`
				ps.id AS submission_id,
				cs.user_id AS submission_user_id,
				COALESCE(p.title, '') AS problem_title,
				ps.state,
				ps.language,
				ps.exec_time AS exec_limit,
				ps.memory_usage AS memory_limit,
				ps.created_at AS submission_time
			`).
			Joins("LEFT JOIN problems AS p ON p.id = cs.problem_id").
			Order("ps.created_at DESC, ps.id DESC").
			Offset(offset).
			Limit(limit).
			Scan(&data.List).Error
	})

	return
}

func getCompetitionSubmissionsBaseQuery(tx *gorm.DB, competitionID int64) *gorm.DB {
	return tx.
		Table("competition_submissions AS cs").
		Joins("JOIN problem_submissions AS ps ON ps.id = cs.submission_id").
		Where("cs.competition_id = ?", competitionID)
}
