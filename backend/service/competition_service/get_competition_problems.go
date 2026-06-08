package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"time"
)

func GetCompetitionProblems(userID int64, req *competition_model.GetCompetitionProblemsReq) (resp response.Response) {

	now := time.Now()
	hasNotStarted, msg := HasNotStarted(req.CompetitionID, now)
	if hasNotStarted {
		resp.Code = 1
		resp.Message = msg
		return
	}

	problems, err := competition_mysql.GetCompetitionProblems(req.CompetitionID, userID)
	if err != nil {
		resp.Code = 1
		resp.Message = "数据查询错误"
		return
	}

	resp.Data = problems
	return
}
