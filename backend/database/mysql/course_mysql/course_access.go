package course_mysql

import (
	"errors"

	"GO1/global"
	"GO1/models/course_model"

	"gorm.io/gorm"
)

type CourseAccessStatus struct {
	CourseID  int64
	Owned     bool
	Creator   bool
	CanAccess bool
}

func GetUserCourseAccessStatus(userID, courseID int64) (status CourseAccessStatus, exist bool, err error) {
	if userID <= 0 || courseID <= 0 {
		return status, false, nil
	}

	var row struct {
		ID      int64
		Owned   int8
		Creator int8
	}
	err = global.DB.
		Table("courses c").
		Select(`
			c.id,
			CASE WHEN c.created_by = ? THEN 1 ELSE 0 END AS creator,
			CASE WHEN EXISTS (
				SELECT 1
				FROM course_enrollments ce
				WHERE ce.course_id = c.id
					AND ce.user_id = ?
					AND ce.status = ?
			) THEN 1 ELSE 0 END AS owned
		`, userID, userID, course_model.CourseEnrollmentStatusActive).
		Where("c.id = ?", courseID).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return status, false, nil
	}
	if err != nil {
		return status, false, err
	}

	status = CourseAccessStatus{
		CourseID:  row.ID,
		Owned:     row.Owned == 1,
		Creator:   row.Creator == 1,
		CanAccess: row.Owned == 1 || row.Creator == 1,
	}
	return status, true, nil
}

func UserCanAccessCourse(userID, courseID int64) (bool, error) {
	status, exist, err := GetUserCourseAccessStatus(userID, courseID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}

	return status.CanAccess, nil
}

func UserCanAccessVideoCourse(userID, videoID int64) (bool, error) {
	if userID <= 0 || videoID <= 0 {
		return false, nil
	}

	courseIDs, err := GetCourseIDsByVideoID(videoID)
	if err != nil {
		return false, err
	}

	for _, courseID := range courseIDs {
		status, exist, err := GetUserCourseAccessStatus(userID, courseID)
		if err != nil {
			return false, err
		}
		if exist && status.CanAccess {
			return true, nil
		}
	}

	return false, nil
}
