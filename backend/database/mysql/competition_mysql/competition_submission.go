package competition_mysql

import (
	"GO1/global"
	"GO1/models/competition_model"
	"GO1/models/problem_model"
	"GO1/models/problem_submission_model"
	"GO1/models/user_ac_problem"
	"GO1/pkg/constants"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveCompetitionJudgeResult(
	competitionID int64,
	problemID int,
	userID int64,
	problemSubmission *problem_submission_model.ProblemSubmission,
	accepted bool,
) (isFirstAC bool, err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		var competitionProblem competition_model.CompetitionProblem
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("competition_id = ? AND problem_id = ?", competitionID, problemID).
			Take(&competitionProblem).Error; err != nil {
			return err
		}

		if accepted {
			var acceptedCount int64
			if err := tx.
				Table("competition_submissions AS cs").
				Joins("JOIN problem_submissions AS ps ON ps.id = cs.submission_id").
				Where("cs.competition_id = ? AND cs.problem_id = ? AND cs.user_id = ? AND ps.state = ?", competitionID, problemID, userID, constants.JudgeStatusAccepted).
				Count(&acceptedCount).Error; err != nil {
				return err
			}
			isFirstAC = acceptedCount == 0
		}

		if err := tx.Create(problemSubmission).Error; err != nil {
			return err
		}

		if err := tx.Create(&competition_model.CompetitionSubmission{
			CompetitionID: competitionID,
			ProblemID:     problemID,
			SubmissionID:  problemSubmission.Id,
			UserID:        userID,
			IsFirstAC:     isFirstAC,
		}).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&problem_model.Problem{}).
			Where("id = ?", problemID).
			UpdateColumn("submit_count", gorm.Expr("submit_count + ?", 1)).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&competition_model.CompetitionProblem{}).
			Where("competition_id = ? AND problem_id = ?", competitionID, problemID).
			UpdateColumn("submit_count", gorm.Expr("submit_count + ?", 1)).Error; err != nil {
			return err
		}

		if !accepted {
			return nil
		}

		if err := tx.
			Clauses(clause.Insert{Modifier: "IGNORE"}).
			Create(&user_ac_problem.UserAcProblem{
				UserId:    userID,
				ProblemId: uint(problemID),
			}).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&problem_model.Problem{}).
			Where("id = ?", problemID).
			UpdateColumn("ac_count", gorm.Expr("ac_count + ?", 1)).Error; err != nil {
			return err
		}

		if !isFirstAC {
			return nil
		}

		return tx.
			Model(&competition_model.CompetitionProblem{}).
			Where("competition_id = ? AND problem_id = ?", competitionID, problemID).
			UpdateColumn("ac_count", gorm.Expr("ac_count + ?", 1)).Error
	})

	return
}
