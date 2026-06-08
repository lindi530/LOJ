package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"
	"GO1/pkg/password_hash"

	"github.com/gin-gonic/gin"
)

func EnterCompetition(req *competition_model.EnterCompetitionReq, c *gin.Context) (resp response.Response) {
	respData := competition_model.EnterCompetitionResp{}
	competition, exist, err := competition_mysql.GetCompetitionAuthInfo(req.CompetitionID)
	if err != nil {
		resp.Code = 1
		resp.ErrCode = response.MysqlError
		return
	}
	if !exist {
		resp.Code = 1
		resp.ErrCode = response.CompetitionNotFound
		return
	}

	if !competition.Visibility && !password_hash.CheckPassword(competition.PasswordHash, req.CompetitionPassword) {
		resp.Code = 1
		resp.Message = "竞赛密码错误"
		return
	}

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
