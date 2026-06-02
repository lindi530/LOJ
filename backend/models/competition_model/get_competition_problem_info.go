package competition_model

import "GO1/models/problem_model"

type GetCompetitionProblemInfoReq struct {
	CompetitionID int64  `json:"competition_id" uri:"competition_id" binding:"required"`
	ProblemNumber string `json:"problem_number" uri:"problem_number" binding:"required"`
}

type GetCompetitionProblemInfoResp struct {
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	InputDesc   string                    `json:"input_desc"`
	OutputDesc  string                    `json:"output_desc"`
	Constraints problem_model.Constraints `json:"constraints" gorm:"-"`
	Language    []string                  `json:"language" gorm:"-"`
	Examples    []problem_model.Example   `json:"examples" gorm:"-"`
}
