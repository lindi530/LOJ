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

func GetVideoProfileByResolution(videoID int64, resolution string) (data course_model.VideoProfile, err error) {
	err = global.DB.
		Table("video_profile").
		Where("video_id = ? AND LOWER(resolution) = LOWER(?)", videoID, resolution).
		First(&data).Error

	return
}

func ChapterVideoExists(courseID, chapterID, videoID int64) (bool, error) {
	var count int64
	err := global.DB.
		Table("chapter_videos").
		Where("course_id = ? AND chapter_id = ? AND video_id = ?", courseID, chapterID, videoID).
		Count(&count).Error

	return count > 0, err
}
