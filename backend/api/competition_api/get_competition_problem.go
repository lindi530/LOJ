package competition_api

import (
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/service/competition_service"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) GetCompetitionProblems(c *gin.Context) {
	req := competition_model.GetCompetitionProblemsReq{}

	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	resp := competition_service.GetCompetitionProblems(&req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
