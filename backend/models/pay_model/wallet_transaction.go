package pay_model

import "time"

type WalletTransaction struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TxNo         string    `gorm:"column:tx_no;size:64;not null;unique" json:"tx_no"`
	UserID       int64     `gorm:"column:user_id;not null;index:idx_user_created_at" json:"user_id"`
	ChangeAmount string    `gorm:"column:change_amount;type:decimal(18,2);not null" json:"change_amount"`
	BalanceAfter string    `gorm:"column:balance_after;type:decimal(18,2);not null" json:"balance_after"`
	BizType      string    `gorm:"column:biz_type;size:32;not null;uniqueIndex:uk_biz_user" json:"biz_type"`
	BizID        string    `gorm:"column:biz_id;size:64;not null;uniqueIndex:uk_biz_user" json:"biz_id"`
	Remark       *string   `gorm:"column:remark;size:255" json:"remark"`
	CreatedAt    time.Time `gorm:"column:created_at;index:idx_user_created_at" json:"created_at"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}
