package user_api

import (
	"GO1/middlewares/response"
	"GO1/models/user_model"
	"GO1/pkg/jwt"
	"GO1/service/user_service"

	"github.com/gin-gonic/gin"
)

func (UserAPI) CheckUserWalletBalance(c *gin.Context) {
	claims, err := jwt.GetUserClaims(c.GetHeader("Authorization"))
	if err != nil || claims.UserId <= 0 {
		response.FailWithCode(response.InvalidAccessToken, c)
		return
	}

	req := user_model.CheckUserWalletBalanceReq{
		UserID: claims.UserId,
	}
	resp := user_service.CheckUserWalletBalance(&req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
