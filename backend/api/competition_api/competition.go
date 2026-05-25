package competition_api

import (
	"GO1/middlewares/response"
	"GO1/service/competition_service"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) GetCompetitions(c *gin.Context) {
	resp := competition_service.GetCompetitions()

	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
