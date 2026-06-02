package competition_model

import "time"

type GetCompetitionSubmissionsReq struct {
	Page          int    `form:"page" json:"page"`
	PageSize      int    `form:"page_size" json:"page_size"`
	ProblemNumber string `form:"problem_number" json:"problem_number"`
}

type GetCompetitionSubmissionsResp struct {
	List           []CompetitionSubmissionsResp `json:"list"`
	Count          int64                        `json:"count"`
	ProblemNumbers []string                     `json:"problem_numbers"`
}

type CompetitionSubmissionsResp struct {
	SubmissionID       int64     `json:"submission_id"`
	SubmissionUserID   int64     `json:"submission_user_id"`
	SubmissionUserName string    `json:"submission_user_name"`
	ProblemTitle       string    `json:"problem_title"`
	ProblemNumber      string    `json:"problem_number"`
	State              string    `json:"state"`
	Language           string    `json:"language"`
	ExecLimit          int       `json:"exec_limit"`
	MemoryLimit        int       `json:"memory_limit"`
	SubmissionTime     time.Time `json:"submission_time"`
}
