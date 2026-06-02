package competition_api

import (
	"GO1/api/problem_api"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"
	"GO1/pkg/jwt"
	"GO1/service/competition_service"

	"github.com/gin-gonic/gin"
)

func (CompetitionAPI) SubmitCompetitionProblem(c *gin.Context) {
	uriReq := struct {
		CompetitionID int64  `uri:"competition_id" binding:"required"`
		ProblemNumber string `uri:"problem_number" binding:"required"`
	}{}

	bodyReq := struct {
		Language string `json:"language" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		global.Logger.Error("submit resp: 1")
		response.FailWithCode(response.BadRequest, c)
		return
	}

	if err := c.ShouldBindJSON(&bodyReq); err != nil {
		global.Logger.Error("submit resp: 2")
		response.FailWithCode(response.BadRequest, c)
		return
	}

	req := competition_model.SubmitCompetitionProblemReq{
		CompetitionID: uriReq.CompetitionID,
		ProblemNumber: uriReq.ProblemNumber,
		Language:      bodyReq.Language,
		Code:          bodyReq.Code,
	}

	if ok := problem_api.IsSafeCode(req.Code, req.Language); !ok {
		response.FailWithMessage(constants.CompetitionSubmitMessageUnsafeCode, c)
		return
	}

	claims, err := jwt.GetUserClaims(c.GetHeader("Authorization"))
	if err != nil {
		response.FailWithCode(response.InvalidAccessToken, c)
		return
	}

	resp := competition_service.SubmitCompetitionProblem(claims.UserId, claims.UserName, &req)

	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
