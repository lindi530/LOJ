package competition_service

import (
	"time"

	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/models/problem_model"
	"GO1/pkg/constants"
)

func GetCompetitionProblemInfo(req *competition_model.GetCompetitionProblemInfoReq) (resp response.Response) {
	now := time.Now()
	hasOver, msg := HasOver(req.CompetitionID, now)
	if hasOver {
		resp.Code = 1
		resp.Message = msg
		return
	}

	problemID, err := getCompetitionProblemID(req.CompetitionID, req.ProblemNumber, now)
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
		return
	}

	problem, constraints, examples, err := competition_mysql.GetCompetitionProblemInfoByID(problemID)
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
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
	resp.Data = problem
	return
}

func problemConstraintToResponse(constraint problem_model.ProblemConstraint) problem_model.Constraint {
	return problem_model.Constraint{
		TimeLimit:   constraint.TimeLimit,
		MemoryLimit: constraint.MemoryLimit,
	}
}
