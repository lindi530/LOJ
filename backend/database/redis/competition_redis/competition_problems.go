package competition_redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"GO1/global"
	"GO1/models/competition_model"

	"github.com/go-redis/redis/v8"
)

func GetCompetitionProblemID(competitionID int64, problemNumber string) (problemID int, exists bool, err error) {
	value, err := global.RedisClient.
		HGet(context.Background(), competitionProblemIDsKey(competitionID), problemNumber).
		Result()
	if errors.Is(err, redis.Nil) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}

	problemID, err = strconv.Atoi(value)
	if err != nil {
		return 0, true, err
	}

	return problemID, true, nil
}

func GetCompetitionProblemNumbers(competitionID int64) (problemNumbers []string, exists bool, err error) {
	problemNumbers, err = global.RedisClient.
		LRange(context.Background(), competitionProblemNumbersKey(competitionID), 0, -1).
		Result()
	if err != nil {
		return nil, false, err
	}
	if len(problemNumbers) == 0 {
		return nil, false, nil
	}

	return problemNumbers, true, nil
}

func GetCompetitionProblems(competitionID int64) (data []competition_model.CompetitionProblem, exists bool, err error) {
	ctx := context.Background()
	problemNumbers, exists, err := GetCompetitionProblemNumbers(competitionID)
	if err != nil || !exists {
		return nil, exists, err
	}

	values, err := global.RedisClient.
		HMGet(ctx, competitionProblemIDsKey(competitionID), problemNumbers...).
		Result()
	if err != nil {
		return nil, false, err
	}

	data = make([]competition_model.CompetitionProblem, 0, len(problemNumbers))
	for i, problemNumber := range problemNumbers {
		if i >= len(values) || values[i] == nil {
			return nil, false, nil
		}

		problemID, err := strconv.Atoi(fmt.Sprint(values[i]))
		if err != nil {
			return nil, true, err
		}

		data = append(data, competition_model.CompetitionProblem{
			CompetitionID: competitionID,
			ProblemNumber: problemNumber,
			ProblemID:     problemID,
		})
	}

	return data, true, nil
}

func SaveCompetitionProblems(
	competitionID int64,
	problems []competition_model.CompetitionProblem,
	endTime time.Time,
) error {
	if len(problems) == 0 || endTime.IsZero() {
		return nil
	}

	expireAt := endTime.Add(time.Hour)
	if !time.Now().Before(expireAt) {
		return nil
	}

	ctx := context.Background()
	idsKey := competitionProblemIDsKey(competitionID)
	numbersKey := competitionProblemNumbersKey(competitionID)
	pipe := global.RedisClient.Pipeline()
	pipe.Del(ctx, idsKey, numbersKey)

	problemIDs := make(map[string]interface{}, len(problems))
	for _, problem := range problems {
		problemIDs[problem.ProblemNumber] = problem.ProblemID
		pipe.RPush(ctx, numbersKey, problem.ProblemNumber)
	}
	pipe.HSet(ctx, idsKey, problemIDs)
	pipe.ExpireAt(ctx, idsKey, expireAt)
	pipe.ExpireAt(ctx, numbersKey, expireAt)

	_, err := pipe.Exec(ctx)
	return err
}

func competitionProblemIDsKey(competitionID int64) string {
	return fmt.Sprintf("competition:problem_ids:%d", competitionID)
}

func competitionProblemNumbersKey(competitionID int64) string {
	return fmt.Sprintf("competition:problem_numbers:%d", competitionID)
}
