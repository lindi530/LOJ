package competition_redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"GO1/global"

	"github.com/go-redis/redis/v8"
)

const competitionRankingScoreBase int64 = 1000000000

func UpdateCompetitionRankingList(
	competitionID int64,
	startTime time.Time,
	endTime time.Time,
	submitTime time.Time,
	problemID int,
	problemNumber string,
	userID int64,
	userName string,
	accepted bool,
) error {
	ctx := context.Background()
	userIDStr := strconv.FormatInt(userID, 10)
	problemIDStr := strconv.Itoa(problemID)
	rankingKey := competitionRankingListKey(competitionID)
	userNameKey := competitionRankingUserNameKey(competitionID)
	userProblemKey := competitionRankingUserProblemKey(competitionID, userID)
	problemNumberField := competitionRankingProblemField(problemIDStr, "problem_number")
	isACField := competitionRankingProblemField(problemIDStr, "is_ac")
	wrongCountField := competitionRankingProblemField(problemIDStr, "wrong_count")
	firstACTimeField := competitionRankingProblemField(problemIDStr, "first_ac_time")

	var expireAt time.Time
	if !endTime.IsZero() {
		expireAt = endTime.Add(time.Hour)
		if !time.Now().Before(expireAt) {
			return nil
		}
	}

	if err := global.RedisClient.ZAddNX(ctx, rankingKey, &redis.Z{
		Score:  0,
		Member: userIDStr,
	}).Err(); err != nil {
		return err
	}

	if err := global.RedisClient.HSet(ctx, userNameKey, userIDStr, userName).Err(); err != nil {
		return err
	}
	if err := global.RedisClient.HSet(ctx, userProblemKey, problemNumberField, problemNumber).Err(); err != nil {
		return err
	}
	if err := global.RedisClient.HSetNX(ctx, userProblemKey, wrongCountField, 0).Err(); err != nil {
		return err
	}

	solvedValue, err := global.RedisClient.HGet(ctx, userProblemKey, isACField).Result()
	if errors.Is(err, redis.Nil) {
		solvedValue = ""
	} else if err != nil {
		return err
	}
	solved := parseRedisBool(solvedValue)

	if !accepted {
		if !solved {
			if err := global.RedisClient.HIncrBy(ctx, userProblemKey, wrongCountField, 1).Err(); err != nil {
				return err
			}
			if err := global.RedisClient.HSet(ctx, userProblemKey, isACField, false).Err(); err != nil {
				return err
			}
		}
		return expireCompetitionRankingKeys(ctx, expireAt, rankingKey, userNameKey, userProblemKey)
	}

	if !solved {
		wrongCount, err := global.RedisClient.HGet(ctx, userProblemKey, wrongCountField).Int64()
		if errors.Is(err, redis.Nil) {
			wrongCount = 0
		} else if err != nil {
			return err
		}

		penalty := competitionRankingPenaltyMinutes(startTime, submitTime) + int(wrongCount)*20
		scoreDelta := competitionRankingScoreBase - int64(penalty)
		if err := global.RedisClient.ZIncrBy(ctx, rankingKey, float64(scoreDelta), userIDStr).Err(); err != nil {
			return err
		}
		if err := global.RedisClient.HSet(ctx, userProblemKey, map[string]interface{}{
			isACField:        true,
			firstACTimeField: submitTime.Format(time.RFC3339Nano),
		}).Err(); err != nil {
			return err
		}
	} else if err := global.RedisClient.HSet(ctx, userProblemKey, isACField, true).Err(); err != nil {
		return err
	}

	return expireCompetitionRankingKeys(ctx, expireAt, rankingKey, userNameKey, userProblemKey)
}

func competitionRankingListKey(competitionID int64) string {
	return fmt.Sprintf("competition:ranking_list:%d", competitionID)
}

func competitionRankingUserNameKey(competitionID int64) string {
	return fmt.Sprintf("competition:ranking_user_name:%d", competitionID)
}

func competitionRankingUserProblemKey(competitionID int64, userID int64) string {
	return fmt.Sprintf("competition:ranking_user_problem:%d:%d", competitionID, userID)
}

func competitionRankingProblemField(problemID string, field string) string {
	return fmt.Sprintf("%s:%s", problemID, field)
}

func competitionRankingPenaltyMinutes(startTime, submitTime time.Time) int {
	if submitTime.Before(startTime) {
		return 0
	}

	return int(submitTime.Sub(startTime).Minutes())
}

func expireCompetitionRankingKeys(ctx context.Context, expireAt time.Time, keys ...string) error {
	if expireAt.IsZero() {
		return nil
	}

	pipe := global.RedisClient.Pipeline()
	for _, key := range keys {
		pipe.ExpireAt(ctx, key, expireAt)
	}
	_, err := pipe.Exec(ctx)
	return err
}
