package user_mysql

import (
	"errors"
	"strconv"
	"time"

	"GO1/global"
	"GO1/models/pay_model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EnsureUserWallet creates a zero-balance wallet when the user has no wallet yet.
func EnsureUserWallet(userID int64) error {
	return EnsureUserWalletWithTx(global.DB, userID)
}

func EnsureUserWalletWithTx(tx *gorm.DB, userID int64) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}
	if tx == nil {
		tx = global.DB
	}

	now := time.Now()
	wallet := pay_model.UserWallet{
		UserID:    userID,
		Balance:   "0.00",
		CreatedAt: now,
		UpdatedAt: now,
	}

	return tx.
		Clauses(clause.Insert{Modifier: "IGNORE"}).
		Create(&wallet).Error
}

func GetUserWalletBalance(userID int64) (float64, error) {
	if userID <= 0 {
		return 0, nil
	}

	var wallet pay_model.UserWallet
	err := global.DB.
		Table("user_wallets").
		Where("user_id = ?", userID).
		First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	balance, err := strconv.ParseFloat(wallet.Balance, 64)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
