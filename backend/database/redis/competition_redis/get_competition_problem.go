package competition_redis

import (
	"GO1/global"
	"GO1/models/competition_model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetCompetitionProblems(competitionID int64) (data []competition_model.GetCompetitionProblemsResp, exists bool, err error) {
	value, err := global.RedisClient.Get(context.Background(), competitionProblemsKey(competitionID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return nil, false, err
	}

	return data, true, nil
}

func SaveCompetitionProblems(competitionID int64, data []competition_model.GetCompetitionProblemsResp, expiration time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return global.RedisClient.
		Set(context.Background(), competitionProblemsKey(competitionID), value, expiration).
		Err()
}

func competitionProblemsKey(competitionID int64) string {
	return fmt.Sprintf("competition:problems:%d", competitionID)
}
