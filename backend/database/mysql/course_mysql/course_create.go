package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"
)

func CreateCourse(course *course_model.Course) error {
	return global.DB.Create(course).Error
}
