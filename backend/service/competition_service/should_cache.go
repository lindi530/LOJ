package competition_service

import "time"

func shouldCacheCompetitionStatus(startTime time.Time, endTime time.Time, now time.Time) bool {
	if startTime.IsZero() || endTime.IsZero() || now.Before(startTime) {
		return false
	}

	return now.Before(endTime.Add(time.Hour))
}

func shouldCacheCompetitionProblems(startTime time.Time, endTime time.Time, now time.Time) bool {
	return shouldCacheCompetitionStatus(startTime, endTime, now)
}
