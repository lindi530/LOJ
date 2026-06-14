package course_redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"GO1/global"

	"github.com/go-redis/redis/v8"
)

const chunkMD5Expiration = time.Hour

func CheckChunkMD5(ctx context.Context, uploadID int64, chunkID int, md5 string) (bool, error) {
	key := chunkMD5Key(uploadID)

	storedMD5, err := global.RedisClient.HGet(ctx, key, strconv.Itoa(chunkID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return storedMD5 == md5, nil
}

func SaveChunkMD5(ctx context.Context, uploadID int64, chunkID int, md5 string) error {
	key := chunkMD5Key(uploadID)
	pipe := global.RedisClient.Pipeline()
	pipe.HSet(ctx, key, strconv.Itoa(chunkID), md5)
	pipe.Expire(ctx, key, chunkMD5Expiration)

	_, err := pipe.Exec(ctx)
	return err
}

func chunkMD5Key(uploadID int64) string {
	return "course:chunk_md5:" + strconv.FormatInt(uploadID, 10)
}
