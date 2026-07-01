package course_mysql

import (
	"errors"
	"time"

	"GO1/global"
	"GO1/models/course_model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrCourseOrderStatusInvalid = errors.New("course order status invalid")

func PayCourseOrder(orderNo, payChannel, transactionID string, paidAt time.Time) (order course_model.CourseOrder, err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table("course_orders").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_no = ?", orderNo).
			First(&order).Error; err != nil {
			return err
		}

		if order.Status != course_model.CourseOrderStatusPaid {
			if order.Status != course_model.CourseOrderStatusPending {
				return ErrCourseOrderStatusInvalid
			}

			if err := tx.
				Table("course_orders").
				Where("id = ?", order.ID).
				Updates(map[string]any{
					"status":         course_model.CourseOrderStatusPaid,
					"pay_channel":    payChannel,
					"transaction_id": transactionID,
					"paid_at":        paidAt,
				}).Error; err != nil {
				return err
			}

			order.Status = course_model.CourseOrderStatusPaid
			order.PayChannel = &payChannel
			order.TransactionID = &transactionID
			order.PaidAt = &paidAt
		}

		created, err := ensureCourseEnrollment(tx, order.UserID, order.CourseID, &order.ID)
		if err != nil {
			return err
		}
		if created {
			return incrementCourseStudentCount(tx, order.CourseID)
		}

		return nil
	})

	return
}
