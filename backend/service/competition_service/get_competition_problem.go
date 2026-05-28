package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/database/redis/competition_redis"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"errors"
	"time"
)

func GetCompetitionProblems(req *competition_model.GetCompetitionProblemsReq) (resp response.Response) {

	now := time.Now()
	hasOver, msg, cacheEndTime, useCache := HasOver(req.CompetitionID, now)
	if hasOver {
		resp.Code = 1
		resp.Message = msg
		return
	}
	var problems []competition_model.GetCompetitionProblemsResp
	cacheHit := false

	if useCache {
		problems, cacheHit, _ = competition_redis.GetCompetitionProblems(req.CompetitionID)
	}

	if !useCache || !cacheHit {
		err := errors.New("")
		problems, err = competition_mysql.GetCompetitionProblems(req.CompetitionID)
		if err != nil {
			resp.Code = 1
			resp.Message = "数据查询错误"
			return
		}
	}

	for i := range problems {
		problems[i].HasAccess = i%2 == 0
	}

	if useCache && !cacheHit {
		_ = competition_redis.SaveCompetitionProblems(req.CompetitionID, problems, cacheEndTime.Sub(now))
	}

	resp.Data = problems
	return
}
