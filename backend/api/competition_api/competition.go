package competition_api

import (
	"GO1/middlewares/response"
	"GO1/service/competition_service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) GetCompetitions(c *gin.Context) {
	hasEnded, err := strconv.ParseBool(c.Param("has_ended"))
	if err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	page := 1
	if hasEnded {
		page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			response.FailWithCode(response.BadRequest, c)
			return
		}
	}

	resp := competition_service.GetCompetitions(hasEnded, page)

	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
