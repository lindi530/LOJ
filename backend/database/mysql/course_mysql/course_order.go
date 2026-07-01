package course_mysql

import (
	"errors"

	"GO1/global"
	"GO1/models/course_model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetCourseByID(courseID int64) (course course_model.Course, exist bool, err error) {
	err = global.DB.
		Table("courses").
		Where("id = ?", courseID).
		First(&course).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return course, false, nil
	}

	return course, err == nil, err
}

func IsCourseCreator(course course_model.Course, userID int64) bool {
	return userID > 0 && course.CreatedBy != nil && *course.CreatedBy == userID
}

func HasActiveCourseEnrollment(userID, courseID int64) (bool, error) {
	if userID <= 0 || courseID <= 0 {
		return false, nil
	}

	var count int64
	err := global.DB.
		Table("course_enrollments").
		Where("user_id = ? AND course_id = ? AND status = ?", userID, courseID, course_model.CourseEnrollmentStatusActive).
		Count(&count).Error

	return count > 0, err
}

func GetOrCreatePendingCourseOrder(userID, courseID int64, amount float64, orderNo string) (order course_model.CourseOrder, reused bool, err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		var course course_model.Course
		if err := tx.
			Table("courses").
			Select("id").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", courseID).
			First(&course).Error; err != nil {
			return err
		}

		err := tx.
			Table("course_orders").
			Where("user_id = ? AND course_id = ? AND status = ?", userID, courseID, course_model.CourseOrderStatusPending).
			Order("id DESC").
			First(&order).Error
		if err == nil {
			reused = true
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		order = course_model.CourseOrder{
			OrderNo:  orderNo,
			UserID:   userID,
			CourseID: courseID,
			Amount:   amount,
			Status:   course_model.CourseOrderStatusPending,
		}
		return tx.Table("course_orders").Create(&order).Error
	})

	return
}

func EnrollFreeCourse(userID, courseID int64) (bool, error) {
	created := false
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		created, err = ensureCourseEnrollment(tx, userID, courseID, nil)
		if err != nil {
			return err
		}
		if created {
			return incrementCourseStudentCount(tx, courseID)
		}
		return nil
	})

	return created, err
}

func ensureCourseEnrollment(tx *gorm.DB, userID, courseID int64, orderID *int64) (bool, error) {
	var enrollment course_model.CourseEnrollment
	err := tx.
		Table("course_enrollments").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&enrollment).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		enrollment = course_model.CourseEnrollment{
			UserID:   userID,
			CourseID: courseID,
			OrderID:  orderID,
			Status:   course_model.CourseEnrollmentStatusActive,
		}
		return true, tx.Table("course_enrollments").Create(&enrollment).Error
	}

	updates := map[string]any{}
	shouldIncrement := false
	if enrollment.Status != course_model.CourseEnrollmentStatusActive {
		updates["status"] = course_model.CourseEnrollmentStatusActive
		shouldIncrement = true
	}
	if orderID != nil {
		updates["order_id"] = *orderID
	}
	if len(updates) == 0 {
		return false, nil
	}

	err = tx.
		Table("course_enrollments").
		Where("id = ?", enrollment.ID).
		Updates(updates).Error
	return shouldIncrement, err
}

func incrementCourseStudentCount(tx *gorm.DB, courseID int64) error {
	return tx.
		Table("courses").
		Where("id = ?", courseID).
		Update("student_count", gorm.Expr("student_count + ?", 1)).Error
}
