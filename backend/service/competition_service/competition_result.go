package competition_service

import (
	"time"

	"GO1/database/mysql/competition_mysql"
	"GO1/database/redis/competition_redis"
	"GO1/global"
	"GO1/models/competition_model"
)

const (
	competitionStatusNotStarted int8 = 1
	competitionStatusRunning    int8 = 2
	competitionStatusEnded      int8 = 3
)

func ScanCompetitionStatus() {
	now := time.Now()
	if err := startDueCompetitions(now); err != nil {
		global.Logger.Error("scan start competitions failed:", err)
	}
	if err := endDueCompetitions(now); err != nil {
		global.Logger.Error("scan end competitions failed:", err)
	}
}

func startDueCompetitions(now time.Time) error {
	competitions, err := competition_mysql.GetNeedStartCompetitions(now)
	if err != nil {
		return err
	}

	for _, competition := range competitions {
		updated, err := competition_mysql.UpdateCompetitionStatus(
			competition.ID,
			[]int8{competitionStatusNotStarted},
			competitionStatusRunning,
		)
		if err != nil {
			global.Logger.Error("start competition failed:", competition.ID, err)
			continue
		}
		if updated {
			if err := competition_redis.SaveCompetitionStatus(
				competition.ID,
				competitionStatusRunning,
				competition.StartTime,
				competition.EndTime,
			); err != nil {
				global.Logger.Error("save competition status failed:", competition.ID, err)
			}
			problems, err := competition_mysql.GetCompetitionProblemMappings(competition.ID)
			if err != nil {
				global.Logger.Error("get competition problems failed:", competition.ID, err)
				continue
			}
			if err := competition_redis.SaveCompetitionProblems(
				competition.ID,
				problems,
				competition.EndTime,
			); err != nil {
				global.Logger.Error("save competition problems failed:", competition.ID, err)
			}
		}
	}

	return nil
}

func endDueCompetitions(now time.Time) error {
	competitions, err := competition_mysql.GetNeedEndCompetitions(now)
	if err != nil {
		return err
	}

	for _, competition := range competitions {
		if err := SaveCompetitionResult(&competition); err != nil {
			global.Logger.Error("save competition result failed:", competition.ID, err)
			continue
		}
		updated, err := competition_mysql.UpdateCompetitionStatus(
			competition.ID,
			[]int8{competitionStatusNotStarted, competitionStatusRunning},
			competitionStatusEnded,
		)
		if err != nil {
			global.Logger.Error("end competition failed:", competition.ID, err)
			continue
		}
		if updated {
			if err := competition_redis.SaveCompetitionStatus(
				competition.ID,
				competitionStatusEnded,
				competition.StartTime,
				competition.EndTime,
			); err != nil {
				global.Logger.Error("save competition status failed:", competition.ID, err)
			}
		}
	}

	return nil
}

func SaveCompetitionResult(competition *competition_model.Competition) error {
	rankingList, err := competition_redis.GetCompetitionRankingList(competition.ID)
	if err != nil {
		return err
	}

	sortAndRankCompetitionRankingList(rankingList)

	userIDs := getRankingUserIDs(rankingList)
	ratingBefore, err := competition_mysql.GetCompetitionUserRatings(userIDs)
	if err != nil {
		return err
	}
	ratingAfter := CalculateRank(rankingList, ratingBefore)

	results := make([]competition_model.CompetitionResult, 0, len(rankingList))
	for _, item := range rankingList {
		results = append(results, competition_model.CompetitionResult{
			CompetitionId:   competition.ID,
			CompetitionName: competition.Name,
			UserId:          item.UserID,
			CompetitionRank: item.CompetitionRank,
			SolvedCount:     item.SolvedCount,
			Penalty:         item.Penalty,
			RatingBefore:    ratingBefore[item.UserID],
			RatingAfter:     ratingAfter[item.UserID],
		})
	}

	return competition_mysql.SaveCompetitionResults(results)
}

func getRankingUserIDs(rankingList []competition_model.CompetitionRankingListItem) []int64 {
	userIDs := make([]int64, 0, len(rankingList))
	for _, item := range rankingList {
		userIDs = append(userIDs, item.UserID)
	}
	return userIDs
}
