package user_model

type CheckUserWalletBalanceReq struct {
	UserID int64 `json:"-"`
}

type CheckUserWalletBalanceResp struct {
	Balance float64 `json:"balance"`
}
