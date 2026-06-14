package course_mysql

import "GO1/global"

func UpdateVideoAssetURL(videoAssetID int64, url string) error {
	return global.DB.Table("video_assets").
		Where("id = ?", videoAssetID).
		Update("url", url).Error
}

func UpdateVideoAssetFinish(videoAssetID int64, url string, sizeBytes int64, duration int) error {
	return global.DB.Table("video_assets").
		Where("id = ?", videoAssetID).
		Updates(map[string]any{
			"url":        url,
			"size_bytes": sizeBytes,
			"duration":   duration,
		}).Error
}
