package problem_mysql

import (
	"GO1/global"
	"GO1/models/problem_model"
	"errors"
)

func GetProblemConstraints(problemId64 int64, language string, constraint *problem_model.ProblemConstraint) error {
	problemId := uint(problemId64)
	if err := global.DB.Where("problem_id = ?	AND language = ?", problemId, language).First(constraint).Error; err != nil {
		return errors.New("查询题目限制失败")
	}
	return nil
}

func GetProblemConstraintList(problemId uint) ([]problem_model.ProblemConstraint, error) {
	var constraints []problem_model.ProblemConstraint
	if err := global.DB.
		Select("time_limit, memory_limit, language").
		Where("problem_id = ?", problemId).
		Find(&constraints).Error; err != nil {
		return nil, errors.New("查询题目限制失败")
	}
	return constraints, nil
}
