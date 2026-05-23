package calendar_redis

import (
	"GO1/global"
	"context"
)

func RemoveACProblem(ctx context.Context, key, member string) error {
	return global.RedisClient.SRem(ctx, key, member).Err()
}
