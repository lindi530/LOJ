package user_service

import (
	"GO1/database/mysql/user_mysql"
	"GO1/middlewares/response"
	"GO1/models/user_model"
)

func CheckUserWalletBalance(req *user_model.CheckUserWalletBalanceReq) (resp response.Response) {
	if req.UserID <= 0 {
		resp.Code = 1
		resp.Message = "invalid user id"
		return
	}

	balance, err := user_mysql.GetUserWalletBalance(req.UserID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query user wallet failed"
		return
	}

	resp.Data = &user_model.CheckUserWalletBalanceResp{
		Balance: balance,
	}
	return
}
