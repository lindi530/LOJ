package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"
)

func GetChapterProblems(courseID, chapterID int64) (data []course_model.ChapterProblemInfo, err error) {
	err = global.DB.
		Table("chapter_problems cp").
		Select("p.id AS id, p.title AS title").
		Joins("JOIN problems p ON p.id = cp.problem_id").
		Where("cp.course_id = ? AND cp.chapter_id = ?", courseID, chapterID).
		Order("cp.sort ASC").
		Order("cp.id ASC").
		Scan(&data).Error

	return
}

func GetChapterVideo(courseID, chapterID int64) (data course_model.ChapterVideoInfo, err error) {
	err = global.DB.
		Table("chapter_videos").
		Select("video_id AS id").
		Where("course_id = ? AND chapter_id = ?", courseID, chapterID).
		Order("id ASC").
		Limit(1).
		Scan(&data).Error

	return
}
