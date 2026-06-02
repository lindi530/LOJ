package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"
	"strings"
	"time"
)

func GetCompetitionSubmissions(competitionID int64, req *competition_model.GetCompetitionSubmissionsReq) (resp response.Response) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	respData, err := competition_mysql.GetCompetitionSubmissions(competitionID, (page-1)*pageSize, pageSize, strings.TrimSpace(req.ProblemNumber))
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
		return
	}

	problemNumbers, err := getCompetitionProblemNumbers(competitionID, time.Now())
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
		return
	}
	respData.ProblemNumbers = problemNumbers

	resp.Data = respData

	return
}
