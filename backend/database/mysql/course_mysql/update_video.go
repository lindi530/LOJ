package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"

	"gorm.io/gorm"
)

func UpdateVideoAssetOriginPath(videoAssetID int64, originPath string) error {
	return global.DB.Table("video_assets").
		Where("id = ?", videoAssetID).
		Updates(map[string]any{
			"origin_path": originPath,
			"status":      course_model.VideoAssetStatusUploading,
		}).Error
}

func UpdateVideoAssetFinish(videoAssetID int64, originPath string, sizeBytes int64, duration int) error {
	return global.DB.Table("video_assets").
		Where("id = ?", videoAssetID).
		Updates(map[string]any{
			"origin_path": originPath,
			"status":      course_model.VideoAssetStatusTranscoding,
			"size_bytes":  sizeBytes,
			"duration":    duration,
		}).Error
}

func UpdateVideoAssetTranscoding(videoAssetID int64) error {
	return global.DB.Table("video_assets").
		Where("id = ?", videoAssetID).
		Update("status", course_model.VideoAssetStatusTranscoding).Error
}

func UpdateVideoAssetTranscodeResult(videoAssetID int64, originPath string, playPath string, coverPath string, profiles []course_model.VideoProfile) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("video_profile").
			Where("video_id = ?", videoAssetID).
			Delete(&course_model.VideoProfile{}).Error; err != nil {
			return err
		}

		for i := range profiles {
			profiles[i].ID = 0
			profiles[i].VideoID = videoAssetID
		}
		if len(profiles) > 0 {
			if err := tx.Table("video_profile").Create(&profiles).Error; err != nil {
				return err
			}
		}

		return tx.Table("video_assets").
			Where("id = ?", videoAssetID).
			Updates(map[string]any{
				"origin_path": originPath,
				"play_path":   playPath,
				"cover_path":  coverPath,
				"status":      course_model.VideoAssetStatusPlayable,
			}).Error
	})
}

func UpdateVideoAssetTranscodeFailed(videoAssetID int64) error {
	return global.DB.Table("video_assets").
		Where("id = ?", videoAssetID).
		Update("status", course_model.VideoAssetStatusTranscodeFailed).Error
}
