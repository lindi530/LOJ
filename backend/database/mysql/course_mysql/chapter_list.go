package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"
)

func GetChapterList(courseID int64) (data []course_model.CourseChapter, err error) {
	err = global.DB.
		Table("course_chapters").
		Where("course_id = ?", courseID).
		Order("sort ASC").
		Order("id ASC").
		Find(&data).Error

	return
}
