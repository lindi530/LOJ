package competition_mysql

import (
	"time"

	"GO1/global"
	"GO1/models/competition_model"
	"GO1/pkg/constants"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type competitionRankingSubmission struct {
	UserID        int64     `gorm:"column:user_id"`
	UserName      string    `gorm:"column:user_name"`
	ProblemID     int       `gorm:"column:problem_id"`
	ProblemNumber string    `gorm:"column:problem_number"`
	State         string    `gorm:"column:state"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

type competitionRankingProblemKey struct {
	UserID    int64
	ProblemID int
}

type competitionRankingProblemStat struct {
	solved     bool
	wrongCount int
}

func GetCompetitionRankingList(competitionID int64, startTime time.Time) (data []competition_model.CompetitionRankingListItem, err error) {
	var submissions []competitionRankingSubmission
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		return tx.
			Table("competition_submissions AS cs").
			Select(`
				cs.user_id,
				COALESCE(cpl.user_name, '') AS user_name,
				cs.problem_id,
				cp.problem_number,
				ps.state,
				ps.created_at
			`).
			Joins("JOIN problem_submissions AS ps ON ps.id = cs.submission_id").
			Joins("JOIN competition_problems AS cp ON cp.competition_id = cs.competition_id AND cp.problem_id = cs.problem_id").
			Joins("LEFT JOIN competition_players AS cpl ON cpl.competition_id = cs.competition_id AND cpl.user_id = cs.user_id").
			Where("cs.competition_id = ?", competitionID).
			Order("ps.created_at ASC").
			Order("cs.id ASC").
			Scan(&submissions).Error
	})
	if err != nil {
		return nil, err
	}

	itemIndexes := make(map[int64]int)
	problemStats := make(map[competitionRankingProblemKey]*competitionRankingProblemStat)
	submissionInfoIndexes := make(map[competitionRankingProblemKey]int)

	for _, submission := range submissions {
		itemIndex, exists := itemIndexes[submission.UserID]
		if !exists {
			data = append(data, competition_model.CompetitionRankingListItem{
				UserID:   submission.UserID,
				UserName: submission.UserName,
			})
			itemIndex = len(data) - 1
			itemIndexes[submission.UserID] = itemIndex
		}

		key := competitionRankingProblemKey{
			UserID:    submission.UserID,
			ProblemID: submission.ProblemID,
		}
		submissionInfoIndex, exists := submissionInfoIndexes[key]
		if !exists {
			data[itemIndex].Submissions = append(data[itemIndex].Submissions, competition_model.CompetitionRankingSubmissionInfo{
				ProblemNumber: submission.ProblemNumber,
			})
			submissionInfoIndex = len(data[itemIndex].Submissions) - 1
			submissionInfoIndexes[key] = submissionInfoIndex
		}

		stat := problemStats[key]
		if stat == nil {
			stat = &competitionRankingProblemStat{}
			problemStats[key] = stat
		}
		if stat.solved {
			continue
		}

		if submission.State == constants.JudgeStatusAccepted {
			stat.solved = true
			data[itemIndex].SolvedCount++
			data[itemIndex].Penalty += competitionPenaltyMinutes(startTime, submission.CreatedAt) + stat.wrongCount*20
			data[itemIndex].Submissions[submissionInfoIndex].IsAC = true
			data[itemIndex].Submissions[submissionInfoIndex].WrongCount = stat.wrongCount
			firstACTime := submission.CreatedAt
			data[itemIndex].Submissions[submissionInfoIndex].FirstACTime = &firstACTime
			continue
		}

		stat.wrongCount++
		data[itemIndex].Submissions[submissionInfoIndex].WrongCount = stat.wrongCount
	}

	return data, nil
}

func competitionPenaltyMinutes(startTime, submitTime time.Time) int {
	if submitTime.Before(startTime) {
		return 0
	}

	return int(submitTime.Sub(startTime).Minutes())
}

func SaveCompetitionResults(results []competition_model.CompetitionResult) error {
	if len(results) == 0 {
		return nil
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// Update before result upsert so retries of the same competition do not recount it.
		if err := UpdateCompetitionSummary(tx, results); err != nil {
			return err
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "competition_id"},
				{Name: "user_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"competition_name",
				"competition_rank",
				"solved_count",
				"penalty",
				"rating_before",
				"rating_after",
			}),
		}).Create(&results).Error; err != nil {
			return err
		}

		for _, result := range results {
			if err := tx.Exec(`
				INSERT INTO saber_stats (user_id, rating)
				VALUES (?, ?)
				ON DUPLICATE KEY UPDATE rating = VALUES(rating), updated_at = NOW()
			`, result.UserId, result.RatingAfter).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
