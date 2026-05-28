package competition_model

import "GO1/models/problem_model"

type SubmitCompetitionProblemReq struct {
	CompetitionID int64  `json:"competition_id" uri:"competition_id" binding:"required"`
	ProblemNumber string `json:"problem_number" uri:"problem_number" binding:"required"`
	Language      string `json:"language" binding:"required"`
	Code          string `json:"code" binding:"required"`
}

type SubmitCompetitionProblemResp struct {
	SubmissionID int64                     `json:"submission_id"`
	State        string                    `json:"state"`
	Score        float64                   `json:"score"`
	IsFirstAC    bool                      `json:"is_first_ac"`
	Results      []problem_model.RunResult `json:"results"`
}
