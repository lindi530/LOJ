package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
)

const endedCompetitionPageSize = 8

func GetCompetitions(hasEnded bool, page int) (resp response.Response) {
	respData := &competition_model.CompetitionResp{}

	competitions, err := competition_mysql.GetCompetitions(hasEnded, page, endedCompetitionPageSize)
	if err != nil {
		resp.Code = 1
		resp.Message = "数据查询错误"
		return
	}

	respData.Competitions = competitions
	resp.Data = respData
	return
}
