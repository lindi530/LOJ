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

func GetCompetitionProblemInfo(competitionID int64, problemNumber string) (data competition_model.GetCompetitionProblemInfoResp, exists bool, err error) {
	value, err := global.RedisClient.Get(context.Background(), competitionProblemInfoKey(competitionID, problemNumber)).Result()
	if errors.Is(err, redis.Nil) {
		return data, false, nil
	}
	if err != nil {
		return data, false, err
	}

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return data, false, err
	}

	return data, true, nil
}

func SaveCompetitionProblemInfo(competitionID int64, problemNumber string, data competition_model.GetCompetitionProblemInfoResp, expiration time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return global.RedisClient.
		Set(context.Background(), competitionProblemInfoKey(competitionID, problemNumber), value, expiration).
		Err()
}

func competitionProblemInfoKey(competitionID int64, problemNumber string) string {
	return fmt.Sprintf("competition:problem_info:%d:%s", competitionID, problemNumber)
}
