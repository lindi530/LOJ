package video

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"GO1/database/mysql/course_mysql"
	"GO1/global"
	"GO1/pkg/constants"
)

func StartVideoTranscodeConsumer() {
	if global.MQChannel == nil {
		global.Logger.Error("video transcode queue unavailable")
		return
	}

	msgs, err := global.MQChannel.Consume(
		constants.VideoTranscodeQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		global.Logger.Error("failed to register video transcode consumer:", err)
		return
	}

	go func() {
		for msg := range msgs {
			var task videoTranscodeTask
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				global.Logger.Error("invalid video transcode message:", err)
				_ = msg.Nack(false, false)
				continue
			}

			if err := processVideoTranscodeTask(context.Background(), task); err != nil {
				global.Logger.Errorw(
					"video transcode task failed",
					"video_asset_id", task.VideoAssetID,
					"bucket", task.Bucket,
					"origin_object", task.OriginObject,
					"err", err,
				)
				if task.VideoAssetID > 0 {
					if updateErr := course_mysql.UpdateVideoAssetTranscodeFailed(task.VideoAssetID); updateErr != nil {
						global.Logger.Errorw(
							"update video asset transcode failed status failed",
							"video_asset_id", task.VideoAssetID,
							"err", updateErr,
						)
					}
				}
				_ = msg.Nack(false, false)
				continue
			}

			if err := msg.Ack(false); err != nil {
				global.Logger.Error("ack video transcode message failed:", err)
			}
		}
	}()
}

func processVideoTranscodeTask(ctx context.Context, task videoTranscodeTask) error {
	if task.VideoAssetID <= 0 || task.Bucket == "" || task.OriginObject == "" {
		return fmt.Errorf("invalid video transcode task")
	}

	if err := course_mysql.UpdateVideoAssetTranscoding(task.VideoAssetID); err != nil {
		return err
	}

	sourcePath, _, _, err := downloadMinioObjectWithMD5(ctx, task.Bucket, task.OriginObject)
	if err != nil {
		return err
	}
	defer os.Remove(sourcePath)

	workDir, err := os.MkdirTemp("", "loj-video-transcode-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workDir)

	hlsDir := filepath.Join(workDir, "hls")
	playlistPath, profiles, err := generateHLS(ctx, sourcePath, hlsDir, task.VideoAssetID)
	if err != nil {
		return err
	}

	coverFilePath := filepath.Join(workDir, "cover.jpg")
	if err := generateCover(ctx, sourcePath, coverFilePath); err != nil {
		return err
	}

	playObjectPrefix := videoAssetObjectPrefix(videoPlayObjectPrefix(), task.VideoAssetID)
	if err := uploadHLSKeys(ctx, task.Bucket, hlsDir, videoHLSKeyObjectPrefix(task.VideoAssetID)); err != nil {
		return err
	}
	if err := uploadHLSFiles(ctx, task.Bucket, hlsDir, playObjectPrefix); err != nil {
		return err
	}

	coverObject := path.Join(videoCoverObjectPrefix(), strconv.FormatInt(task.VideoAssetID, 10)+".jpg")
	if err := uploadLocalFileToMinio(ctx, task.Bucket, coverObject, coverFilePath, "image/jpeg"); err != nil {
		return err
	}

	playObject := path.Join(playObjectPrefix, path.Base(filepath.ToSlash(playlistPath)))
	originPath := task.OriginPath
	if originPath == "" {
		originPath = minioVideoPath(task.Bucket, task.OriginObject)
	}

	for i := range profiles {
		profiles[i].VideoID = task.VideoAssetID
		profiles[i].PlaylistPath = minioVideoPath(
			task.Bucket,
			path.Join(playObjectPrefix, profiles[i].PlaylistPath),
		)
	}

	return course_mysql.UpdateVideoAssetTranscodeResult(
		task.VideoAssetID,
		originPath,
		minioVideoPath(task.Bucket, playObject),
		minioVideoPath(task.Bucket, coverObject),
		profiles,
	)
}
