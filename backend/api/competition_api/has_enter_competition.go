package competition_api

import (
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/jwt"
	"GO1/service/competition_service"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) HasEnterCompetition(c *gin.Context) {
	req := competition_model.HasEnterCompetitionReq{}

	jwt.SaveUserInfoFromToken(c)

	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	resp := competition_service.HasEnterCompetition(c, &req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
