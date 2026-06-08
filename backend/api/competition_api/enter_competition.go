package competition_api

import (
	"errors"
	"io"

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
	competitionID := req.CompetitionID
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.FailWithCode(response.ReqError, c)
		return
	}
	req.CompetitionID = competitionID

	resp := competition_service.EnterCompetition(&req, c)

	if resp.Code == 1 {
		if resp.Message != nil {
			response.FailWithMessage(resp.Message, c)
			return
		}
		response.FailWithCode(resp.ErrCode, c)
		return
	}

	response.OkWithData(resp, c)
}
