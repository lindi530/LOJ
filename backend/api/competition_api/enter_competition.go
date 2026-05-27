package competition_api

import (
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/jwt"
	"GO1/service/competition_service"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) EnterCompetition(c *gin.Context) {
	jwt.SaveUserInfoFromToken(c)

	req := competition_model.EnterCompetitionReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.ReqError, c)
		return
	}

	resp := competition_service.EnterCompetition(&req, c)

	if resp.Code == 1 {
		response.FailWithCode(resp.ErrCode, c)
		return
	}

	response.OkWithData(resp, c)
}
