package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/database/redis/competition_redis"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/models/problem_model"
	"time"
)

func GetCompetitionProblemInfo(req *competition_model.GetCompetitionProblemInfoReq) (resp response.Response) {
	now := time.Now()
	hasOver, msg, cacheEndTime, useCache := HasOver(req.CompetitionID, now)
	if hasOver {
		resp.Code = 1
		resp.Message = msg
		return
	}

	if useCache {
		problem, cacheHit, cacheErr := competition_redis.GetCompetitionProblemInfo(req.CompetitionID, req.ProblemNumber)
		if cacheErr == nil && cacheHit {
			resp.Data = problem
			return
		}
	}

	problem, constraints, examples, err := competition_mysql.GetCompetitionProblemInfo(req.CompetitionID, req.ProblemNumber)
	if err != nil {
		resp.Code = 1
		resp.Message = "数据查询错误"
		return
	}

	for _, constraint := range constraints {
		problem.Language = append(problem.Language, constraint.Language)
		switch constraint.Language {
		case "cpp":
			problem.Constraints.Cpp = problemConstraintToResponse(constraint)
		case "python":
			problem.Constraints.Python = problemConstraintToResponse(constraint)
		case "go":
			problem.Constraints.Go = problemConstraintToResponse(constraint)
		}
	}

	problem.Examples = examples

	if useCache {
		_ = competition_redis.SaveCompetitionProblemInfo(
			req.CompetitionID,
			req.ProblemNumber,
			problem,
			cacheEndTime.Sub(now),
		)
	}

	resp.Data = problem
	return
}

func problemConstraintToResponse(constraint problem_model.ProblemConstraint) problem_model.Constraint {
	return problem_model.Constraint{
		TimeLimit:   constraint.TimeLimit,
		MemoryLimit: constraint.MemoryLimit,
	}
}
