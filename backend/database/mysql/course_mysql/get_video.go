package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"
)

func GetVideoAssetByID(videoAssetID int64) (data course_model.VideoAsset, err error) {
	err = global.DB.Table("video_assets").
		Where("id = ?", videoAssetID).
		First(&data).Error
	return
}
