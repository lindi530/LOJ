package problem_service

import (
	mysql "GO1/database/mysql/problem_mysql"
	"GO1/middlewares/response"
	models "GO1/models/problem_model"
)

func GetProblemDetails(id uint) (resp response.Response) {
	var problemDetail models.Problem

	if err := mysql.GetProblemDetails(id, &problemDetail); err != nil {
		resp.Code = 1
		resp.Message = "读取题目信息失败"
		return
	}

	if err := mysql.GetProblemTags(id, &problemDetail.Tags); err != nil {
		resp.Code = 1
		resp.Message = "读取题目信息失败"
		return
	}

	if err := mysql.GetOnProblemExamples(id, &problemDetail.Examples); err != nil {
		resp.Code = 1
		resp.Message = "读取题目信息失败"
		return
	}

	constraints, err := mysql.GetProblemConstraintList(id)
	if err != nil {
		resp.Code = 1
		resp.Message = "读取题目信息失败"
		return
	}
	for _, constraint := range constraints {
		problemDetail.Language = append(problemDetail.Language, constraint.Language)
		switch constraint.Language {
		case "cpp":
			problemDetail.Constraints.Cpp = problemConstraintToResponse(constraint)
		case "python":
			problemDetail.Constraints.Python = problemConstraintToResponse(constraint)
		case "go":
			problemDetail.Constraints.Go = problemConstraintToResponse(constraint)
		}
	}

	resp.Code = 0
	resp.Data = problemDetail
	resp.Message = "读取题目信息成功"

	return
}

func problemConstraintToResponse(constraint models.ProblemConstraint) models.Constraint {
	return models.Constraint{
		TimeLimit:   constraint.TimeLimit,
		MemoryLimit: constraint.MemoryLimit,
	}
}
