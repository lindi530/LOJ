package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"time"
)

func HasOver(competitionID int64, now time.Time) (bool, string, time.Time, bool) {
	competition, err := competition_mysql.GetCompetitionTime(competitionID)
	if err != nil {
		return true, "竞赛不存在", time.Time{}, false
	}

	if competition.StartTime.IsZero() || now.Before(competition.StartTime) {
		return true, "竞赛未开始", time.Time{}, false
	}

	cacheEndTime := competition.EndTime.Add(time.Hour)
	useCache := now.Before(cacheEndTime)

	return false, "", cacheEndTime, useCache
}
