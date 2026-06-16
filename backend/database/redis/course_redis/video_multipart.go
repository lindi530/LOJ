package course_redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"GO1/global"
)

const videoMultipartExpiration = 24 * time.Hour

type VideoMultipartUploadMeta struct {
	MinioUploadID string `json:"minio_upload_id"`
	Bucket        string `json:"bucket"`
	TempObject    string `json:"temp_object"`
}

func SaveVideoMultipartUpload(ctx context.Context, uploadID int64, meta VideoMultipartUploadMeta) error {
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	return global.RedisClient.Set(ctx, videoMultipartUploadKey(uploadID), data, videoMultipartExpiration).Err()
}

func GetVideoMultipartUpload(ctx context.Context, uploadID int64) (VideoMultipartUploadMeta, error) {
	data, err := global.RedisClient.Get(ctx, videoMultipartUploadKey(uploadID)).Bytes()
	if err != nil {
		return VideoMultipartUploadMeta{}, err
	}

	var meta VideoMultipartUploadMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return VideoMultipartUploadMeta{}, err
	}

	return meta, nil
}

func DeleteVideoMultipartUpload(ctx context.Context, uploadID int64) error {
	return global.RedisClient.Del(ctx, videoMultipartUploadKey(uploadID)).Err()
}

func videoMultipartUploadKey(uploadID int64) string {
	return "course:video_multipart:" + strconv.FormatInt(uploadID, 10)
}
