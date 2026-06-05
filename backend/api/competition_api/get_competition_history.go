package competition_api

import (
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/service/competition_service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) GetCompetitionHistory(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	req := competition_model.GetCompetitionHistoryReq{
		UserID: userID,
	}

	resp := competition_service.GetCompetitionHistory(&req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
