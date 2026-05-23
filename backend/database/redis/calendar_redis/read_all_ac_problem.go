package calendar_redis

import (
	"GO1/global"
	"context"
)

func ReadAllACProblem(ctx context.Context, key string, cursor uint64, count int64) ([]string, uint64, error) {

	return global.RedisClient.SScan(
		ctx,
		key,
		cursor,
		"*",
		count,
	).Result()
}
