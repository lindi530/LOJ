package course_mysql

import "GO1/global"

func GetCourseIDsByVideoID(videoID int64) ([]int64, error) {
	if videoID <= 0 {
		return nil, nil
	}

	var courseIDs []int64
	err := global.DB.
		Table("chapter_videos").
		Distinct("course_id").
		Where("video_id = ?", videoID).
		Pluck("course_id", &courseIDs).Error
	return courseIDs, err
}
