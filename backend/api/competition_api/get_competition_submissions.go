package competition_api

import (
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/jwt"
	"GO1/service/competition_service"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) GetCompetitionSubmissions(c *gin.Context) {
	jwt.SaveUserInfoFromToken(c)

	uriReq := struct {
		CompetitionID int64 `uri:"competition_id" binding:"required"`
	}{}
	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	req := competition_model.GetCompetitionSubmissionsReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithCode(response.ReqError, c)
		return
	}

	resp := competition_service.GetCompetitionSubmissions(uriReq.CompetitionID, &req)

	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
