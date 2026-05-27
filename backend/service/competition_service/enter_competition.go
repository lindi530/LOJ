package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"

	"github.com/gin-gonic/gin"
)

func EnterCompetition(req *competition_model.EnterCompetitionReq, c *gin.Context) (resp response.Response) {
	respData := competition_model.EnterCompetitionResp{}
	enterPlayer := competition_model.CompetitionPlayer{
		CompetitionID: req.CompetitionID,
		UserID:        c.GetInt64(constants.LoginUserIDKey),
		UserName:      c.GetString(constants.LoginUserNameKey),
	}
	if exist, err := competition_mysql.EnterCompetition(&enterPlayer); err != nil {
		resp.Code = 1
		resp.ErrCode = response.MysqlError
		return
	} else if !exist {
		resp.Code = 1
		resp.ErrCode = response.CompetitionNotFound
		return
	}

	resp.Message = "success"
	resp.Data = &respData
	return
}
