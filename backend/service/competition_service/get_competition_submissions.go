package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
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

	respData, err := competition_mysql.GetCompetitionSubmissions(competitionID, (page-1)*pageSize, pageSize)
	if err != nil {
		resp.Code = 1
		resp.Message = "data query error"
		return
	}

	resp.Data = respData

	return
}
