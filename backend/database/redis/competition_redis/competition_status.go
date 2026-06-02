package competition_redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"GO1/global"
)

type CompetitionStatusCache struct {
	Status    int8
	StartTime time.Time
	EndTime   time.Time
}

func GetCompetitionStatus(competitionID int64) (data CompetitionStatusCache, exists bool, err error) {
	ctx := context.Background()
	fields, err := global.RedisClient.HGetAll(ctx, competitionStatusKey(competitionID)).Result()
	if err != nil {
		return data, false, err
	}
	if len(fields) == 0 {
		return data, false, nil
	}

	status, err := strconv.Atoi(fields["status"])
	if err != nil {
		return data, true, err
	}

	startTime, err := parseCompetitionStatusTime(fields["start_time"])
	if err != nil {
		return data, true, err
	}

	endTime, err := parseCompetitionStatusTime(fields["end_time"])
	if err != nil {
		return data, true, err
	}

	data.Status = int8(status)
	data.StartTime = startTime
	data.EndTime = endTime
	return data, true, nil
}

func SaveCompetitionStatus(competitionID int64, status int8, startTime time.Time, endTime time.Time) error {
	if endTime.IsZero() {
		return nil
	}

	expireAt := endTime.Add(time.Hour)
	if !time.Now().Before(expireAt) {
		return nil
	}

	ctx := context.Background()
	key := competitionStatusKey(competitionID)
	pipe := global.RedisClient.Pipeline()
	pipe.HSet(ctx, key, map[string]interface{}{
		"status":     status,
		"start_time": startTime.Format(time.RFC3339Nano),
		"end_time":   endTime.Format(time.RFC3339Nano),
	})
	pipe.ExpireAt(ctx, key, expireAt)
	_, err := pipe.Exec(ctx)
	return err
}

func competitionStatusKey(competitionID int64) string {
	return fmt.Sprintf("competition:status:%d", competitionID)
}

func parseCompetitionStatusTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339Nano, value)
}
