package calendar_redis

import (
	"GO1/global"
	"context"
	"fmt"
	"time"
)

func SaveACProblem(userID int64, problemID uint) error {
	ctx := context.Background()

	date := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("pending:user:problem:%s", date)

	member := fmt.Sprintf("%d:%d", userID, problemID)

	return global.RedisClient.SAdd(ctx, key, member).Err()
}
