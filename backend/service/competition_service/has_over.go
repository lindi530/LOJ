package competition_service

import (
	"time"

	"GO1/database/mysql/competition_mysql"
	"GO1/database/redis/competition_redis"
	"GO1/global"
	"GO1/models/competition_model"
	"GO1/pkg/constants"
)

func HasOver(competitionID int64, now time.Time) (bool, string) {
	cache, exists, err := competition_redis.GetCompetitionStatus(competitionID)
	if err != nil {
		global.Logger.Error("get competition status cache failed:", competitionID, err)
	} else if exists {
		status := normalizeCompetitionStatus(cache.Status, cache.StartTime, cache.EndTime, now)
		if status != cache.Status {
			if err := competition_redis.SaveCompetitionStatus(
				competitionID,
				status,
				cache.StartTime,
				cache.EndTime,
			); err != nil {
				global.Logger.Error("save competition status failed:", competitionID, err)
			}
		}
		return competitionStatusResult(status)
	}

	competition, err := competition_mysql.GetCompetitionTime(competitionID)
	if err != nil {
		return true, constants.CompetitionSubmitMessageDataQueryError
	}

	status := currentCompetitionStatus(competition, now)
	if shouldCacheCompetitionStatus(competition.StartTime, competition.EndTime, now) {
		if err := competition_redis.SaveCompetitionStatus(
			competition.ID,
			status,
			competition.StartTime,
			competition.EndTime,
		); err != nil {
			global.Logger.Error("save competition status failed:", competition.ID, err)
		}
	}

	return competitionStatusResult(status)
}

func currentCompetitionStatus(competition competition_model.Competition, now time.Time) int8 {
	return normalizeCompetitionStatus(
		competition.Status,
		competition.StartTime,
		competition.EndTime,
		now,
	)
}

func normalizeCompetitionStatus(status int8, startTime time.Time, endTime time.Time, now time.Time) int8 {
	if status == competitionStatusEnded || (!endTime.IsZero() && !now.Before(endTime)) {
		return competitionStatusEnded
	}
	if startTime.IsZero() || now.Before(startTime) {
		return competitionStatusNotStarted
	}

	return competitionStatusRunning
}

func competitionStatusResult(status int8) (bool, string) {
	if status == competitionStatusEnded {
		return true, constants.CompetitionSubmitMessageEnded
	}
	if status == competitionStatusNotStarted {
		return true, constants.CompetitionSubmitMessageNotStarted
	}

	return false, ""
}
