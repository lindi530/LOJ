package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"
)

func CreateEmptyVideoAsset() (*course_model.VideoAsset, error) {
	videoAsset := &course_model.VideoAsset{
		OriginPath: "",
		Status:     course_model.VideoAssetStatusUploading,
		Duration:   0,
	}

	err := global.DB.Table("video_assets").Create(videoAsset).Error
	return videoAsset, err
}
