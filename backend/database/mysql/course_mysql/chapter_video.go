package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"
)

func GetVideoProfilesByVideoID(videoID int64) (data []course_model.VideoProfile, err error) {
	err = global.DB.
		Table("video_profile").
		Where("video_id = ?", videoID).
		Order("height ASC").
		Order("id ASC").
		Find(&data).Error

	return
}
