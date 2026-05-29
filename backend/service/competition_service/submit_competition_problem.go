package competition_service

import (
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"
)

type competitionSubmitJob struct {
	UserID        int64  `json:"user_id"`
	CompetitionID int64  `json:"competition_id"`
	ProblemNumber string `json:"problem_number"`
	Language      string `json:"language"`
	Code          string `json:"code"`
}

type competitionSubmitResult struct {
	Code    int                                            `json:"code"`
	Message string                                         `json:"message"`
	Data    competition_model.SubmitCompetitionProblemResp `json:"data"`
}

func SubmitCompetitionProblem(userID int64, req *competition_model.SubmitCompetitionProblemReq) (resp response.Response) {
	if global.MQChannel == nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageJudgeQueueUnavailable
		return
	}

	result, err := publishCompetitionSubmitJob(competitionSubmitJob{
		UserID:        userID,
		CompetitionID: req.CompetitionID,
		ProblemNumber: req.ProblemNumber,
		Language:      req.Language,
		Code:          req.Code,
	})
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		return
	}

	resp.Code = result.Code
	resp.Message = result.Message
	resp.Data = result.Data

	return
}
