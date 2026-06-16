package video

import (
	"path"
	"strconv"
	"strings"

	"GO1/global"
)

type videoTranscodeTask struct {
	VideoAssetID int64  `json:"video_asset_id"`
	Bucket       string `json:"bucket"`
	OriginObject string `json:"origin_object"`
	OriginPath   string `json:"origin_path"`
}

func minioVideoPath(bucket string, object string) string {
	return path.Join(bucket, object)
}

func videoSourceObjectPrefix() string {
	prefix := cleanObjectPrefix(global.Config.Minio.Bucket.VideoPath)
	if prefix == "" {
		return "raw"
	}
	return prefix
}

func videoPlayObjectPrefix() string {
	prefix := cleanObjectPrefix(global.Config.Minio.Bucket.VideoPlayPath)
	if prefix == "" {
		return "hls"
	}
	return prefix
}

func videoCoverObjectPrefix() string {
	prefix := cleanObjectPrefix(global.Config.Minio.Bucket.VideoCoverPath)
	if prefix == "" {
		return "cover"
	}
	return prefix
}

func videoAssetObjectPrefix(prefix string, videoAssetID int64) string {
	return path.Join(prefix, strconv.FormatInt(videoAssetID, 10))
}

func cleanObjectPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	prefix = strings.ReplaceAll(prefix, "\\", "/")
	return strings.Trim(prefix, "/")
}
