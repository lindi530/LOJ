package pay_model

import "time"

type UserWallet struct {
	UserID    int64     `gorm:"column:user_id;primaryKey" json:"user_id"`
	Balance   string    `gorm:"column:balance;type:decimal(18,2);not null;default:0.00" json:"balance"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (UserWallet) TableName() string {
	return "user_wallets"
}
