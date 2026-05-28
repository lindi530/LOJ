package competition_api

import (
	"GO1/api/problem_api"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
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
		response.FailWithMessage("unsafe code", c)
		return
	}

	userID := jwt.GetUserIdFromToken(c.GetHeader("Authorization"))
	resp := competition_service.SubmitCompetitionProblem(userID, &req)

	global.Logger.Error("submit respdata: ", resp)

	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.OkWithData(resp.Data, c)
}
