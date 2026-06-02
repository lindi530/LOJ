package competition_service

import (
	"errors"
	"strings"
	"time"

	"GO1/database/mysql/competition_mysql"
	"GO1/database/redis/competition_redis"
	"GO1/global"
	"GO1/models/competition_model"
)

var errCompetitionProblemNotFound = errors.New("competition problem not found")

func getCompetitionProblemID(competitionID int64, problemNumber string, now time.Time) (int, error) {
	problemNumber = strings.TrimSpace(problemNumber)
	if problemNumber == "" {
		return 0, errCompetitionProblemNotFound
	}

	problemID, exists, err := competition_redis.GetCompetitionProblemID(competitionID, problemNumber)
	if err != nil {
		global.Logger.Error("get competition problem id cache failed:", competitionID, problemNumber, err)
	} else if exists {
		return problemID, nil
	}

	problems, err := loadCompetitionProblemsFromMysqlAndCache(competitionID, now)
	if err != nil {
		return 0, err
	}

	for _, problem := range problems {
		if problem.ProblemNumber == problemNumber {
			return problem.ProblemID, nil
		}
	}

	return 0, errCompetitionProblemNotFound
}

func getCompetitionProblemNumbers(competitionID int64, now time.Time) ([]string, error) {
	problemNumbers, exists, err := competition_redis.GetCompetitionProblemNumbers(competitionID)
	if err != nil {
		global.Logger.Error("get competition problem numbers cache failed:", competitionID, err)
	} else if exists {
		return problemNumbers, nil
	}

	problems, err := loadCompetitionProblemsFromMysqlAndCache(competitionID, now)
	if err != nil {
		return nil, err
	}

	problemNumbers = make([]string, 0, len(problems))
	for _, problem := range problems {
		problemNumbers = append(problemNumbers, problem.ProblemNumber)
	}

	return problemNumbers, nil
}

func loadCompetitionProblemsFromMysqlAndCache(competitionID int64, now time.Time) ([]competition_model.CompetitionProblem, error) {
	problems, err := competition_mysql.GetCompetitionProblemMappings(competitionID)
	if err != nil {
		return nil, err
	}

	cacheCompetitionProblemsIfNeeded(competitionID, problems, now)
	return problems, nil
}

func cacheCompetitionProblemsIfNeeded(competitionID int64, problems []competition_model.CompetitionProblem, now time.Time) {
	if len(problems) == 0 {
		return
	}

	cache, exists, err := competition_redis.GetCompetitionStatus(competitionID)
	if err != nil {
		global.Logger.Error("get competition status cache for problem cache failed:", competitionID, err)
	}

	var startTime, endTime time.Time
	if exists && err == nil {
		startTime = cache.StartTime
		endTime = cache.EndTime
	} else {
		competition, err := competition_mysql.GetCompetitionTime(competitionID)
		if err != nil {
			global.Logger.Error("get competition time for problem cache failed:", competitionID, err)
			return
		}
		startTime = competition.StartTime
		endTime = competition.EndTime
	}

	if !shouldCacheCompetitionProblems(startTime, endTime, now) {
		return
	}

	if err := competition_redis.SaveCompetitionProblems(
		competitionID,
		problems,
		endTime,
	); err != nil {
		global.Logger.Error("save competition problems cache failed:", competitionID, err)
	}
}
