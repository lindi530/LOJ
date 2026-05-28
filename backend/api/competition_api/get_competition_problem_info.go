package competition_api

import (
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/service/competition_service"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) GetCompetitionProblemInfo(c *gin.Context) {
	req := competition_model.GetCompetitionProblemInfoReq{}

	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	resp := competition_service.GetCompetitionProblemInfo(&req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}
	global.Logger.Error("GetCompetitionProblemInfo:", resp.Data)
	response.OkWithData(resp.Data, c)
}
