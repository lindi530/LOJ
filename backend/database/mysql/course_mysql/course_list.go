package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"
)

func GetCourseList() (data []course_model.Course, total int64, err error) {
	if err = global.DB.Table("courses").
		Where("status = ?", 0).
		Count(&total).Error; err != nil {
		return
	}

	err = global.DB.Table("courses").
		Where("status = ?", 0).
		Order("sort ASC").
		Order("id DESC").
		Find(&data).Error

	return
}
