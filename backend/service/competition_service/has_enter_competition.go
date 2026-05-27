package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"

	"github.com/gin-gonic/gin"
)

func HasEnterCompetition(c *gin.Context, req *competition_model.HasEnterCompetitionReq) (resp response.Response) {
	respData := &competition_model.HasEnterCompetitionResp{}

	userID, exists := c.Get(constants.LoginUserIDKey)
	if !exists {
		resp.Code = 1
		resp.Message = "用户ID错误"
		return
	}

	data, err := competition_mysql.HasEnterCompetition(req.CompetitionID, userID.(int64))
	if err != nil {
		resp.Code = 1
		resp.Message = "数据查询错误"
		return
	}

	respData.HasEnter = len(data) > 0
	resp.Data = respData

	return
}
