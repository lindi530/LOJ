package order_mysql

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"GO1/database/mysql/user_mysql"
	"GO1/global"
	"GO1/models/course_model"
	"GO1/models/pay_model"
	pkg_snowflake "GO1/pkg/snowflake"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrCourseOrderNotFound      = errors.New("course order not found")
	ErrCourseOrderUserInvalid   = errors.New("course order user invalid")
	ErrCourseOrderStatusInvalid = errors.New("course order status invalid")
	ErrUserWalletNotFound       = errors.New("user wallet not found")
	ErrCoinBalanceInsufficient  = errors.New("coin balance insufficient")
	ErrUserWalletBalanceInvalid = errors.New("user wallet balance invalid")
	ErrCourseOrderAmountInvalid = errors.New("course order amount invalid")
)

func PayCourseOrderByCoin(userID int64, orderNo string) (order course_model.CourseOrder, balance float64, owned bool, err error) {
	var businessErr error

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table("course_orders").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_no = ?", orderNo).
			First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCourseOrderNotFound
			}
			return err
		}

		if order.UserID != userID {
			return ErrCourseOrderUserInvalid
		}
		if order.Status != course_model.CourseOrderStatusPending {
			return ErrCourseOrderStatusInvalid
		}

		isCreator, err := userIsCourseCreator(tx, userID, order.CourseID)
		if err != nil {
			return err
		}
		if isCreator {
			balance, _ = getUserWalletBalance(tx, userID)

			payChannel := "creator"
			paidAt := time.Now()
			if err := tx.
				Table("course_orders").
				Where("id = ?", order.ID).
				Updates(map[string]any{
					"amount":      0,
					"status":      course_model.CourseOrderStatusPaid,
					"pay_channel": payChannel,
					"paid_at":     paidAt,
				}).Error; err != nil {
				return err
			}

			order.Amount = 0
			order.Status = course_model.CourseOrderStatusPaid
			order.PayChannel = &payChannel
			order.PaidAt = &paidAt
			owned = true
			return nil
		}

		var wallet pay_model.UserWallet
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&wallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := user_mysql.EnsureUserWalletWithTx(tx, userID); err != nil {
					return err
				}
				balance = 0
				businessErr = ErrCoinBalanceInsufficient
				return nil
			}
			return err
		}

		balanceCents, err := decimalStringToCents(wallet.Balance)
		if err != nil {
			return ErrUserWalletBalanceInvalid
		}
		amountCents, err := amountToCents(order.Amount)
		if err != nil {
			return ErrCourseOrderAmountInvalid
		}
		if balanceCents < amountCents {
			return ErrCoinBalanceInsufficient
		}

		balanceAfterCents := balanceCents - amountCents
		balanceAfter := centsToDecimalString(balanceAfterCents)
		if err := tx.
			Table("user_wallets").
			Where("user_id = ?", userID).
			Update("balance", balanceAfter).Error; err != nil {
			return err
		}

		changeAmount := centsToDecimalString(-amountCents)
		remark := "course order coin pay"
		walletTx := pay_model.WalletTransaction{
			TxNo:         fmt.Sprintf("WALLET%d", pkg_snowflake.Snowflake{}.GenID()),
			UserID:       userID,
			ChangeAmount: changeAmount,
			BalanceAfter: balanceAfter,
			BizType:      "course_order",
			BizID:        order.OrderNo,
			Remark:       &remark,
		}
		if err := tx.Table("wallet_transactions").Create(&walletTx).Error; err != nil {
			return err
		}

		payChannel := "coin"
		paidAt := time.Now()
		if err := tx.
			Table("course_orders").
			Where("id = ?", order.ID).
			Updates(map[string]any{
				"status":      course_model.CourseOrderStatusPaid,
				"pay_channel": payChannel,
				"paid_at":     paidAt,
			}).Error; err != nil {
			return err
		}

		order.Status = course_model.CourseOrderStatusPaid
		order.PayChannel = &payChannel
		order.PaidAt = &paidAt

		created, err := ensureCourseEnrollment(tx, userID, order.CourseID, &order.ID)
		if err != nil {
			return err
		}
		if created {
			if err := incrementCourseStudentCount(tx, order.CourseID); err != nil {
				return err
			}
		}

		balance = centsToFloat64(balanceAfterCents)
		owned = true
		return nil
	})
	if err == nil && businessErr != nil {
		err = businessErr
	}

	return
}

func userIsCourseCreator(tx *gorm.DB, userID, courseID int64) (bool, error) {
	if userID <= 0 || courseID <= 0 {
		return false, nil
	}

	var course course_model.Course
	err := tx.
		Table("courses").
		Select("id").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND created_by = ?", courseID, userID).
		First(&course).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return err == nil, err
}

func getUserWalletBalance(tx *gorm.DB, userID int64) (float64, error) {
	if userID <= 0 {
		return 0, nil
	}

	var wallet pay_model.UserWallet
	err := tx.
		Table("user_wallets").
		Where("user_id = ?", userID).
		First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	balanceCents, err := decimalStringToCents(wallet.Balance)
	if err != nil {
		return 0, err
	}
	return centsToFloat64(balanceCents), nil
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

func amountToCents(amount float64) (int64, error) {
	if amount < 0 {
		return 0, errors.New("negative amount")
	}
	return int64(math.Round(amount * 100)), nil
}

func decimalStringToCents(value string) (int64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, errors.New("empty decimal")
	}

	sign := int64(1)
	if value[0] == '-' {
		sign = -1
		value = value[1:]
	}

	parts := strings.Split(value, ".")
	if len(parts) > 2 || parts[0] == "" {
		return 0, errors.New("invalid decimal")
	}

	integerPart, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}

	fraction := "00"
	if len(parts) == 2 {
		fraction = parts[1]
		if len(fraction) > 2 {
			return 0, errors.New("invalid decimal scale")
		}
		fraction += strings.Repeat("0", 2-len(fraction))
	}
	fractionPart, err := strconv.ParseInt(fraction, 10, 64)
	if err != nil {
		return 0, err
	}

	return sign * (integerPart*100 + fractionPart), nil
}

func centsToDecimalString(cents int64) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
}

func centsToFloat64(cents int64) float64 {
	return float64(cents) / 100
}
