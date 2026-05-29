package competition_model

import "time"

type GetCompetitionSubmissionsReq struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

type GetCompetitionSubmissionsResp struct {
	List  []CompetitionSubmissionsResp `json:"list"`
	Count int64                        `json:"count"`
}

type CompetitionSubmissionsResp struct {
	SubmissionID     int64     `json:"submission_id"`
	SubmissionUserID int64     `json:"submission_user_id"`
	ProblemTitle     string    `json:"problem_title"`
	State            string    `json:"state"`
	Language         string    `json:"language"`
	ExecLimit        int       `json:"exec_limit"`
	MemoryLimit      int       `json:"memory_limit"`
	SubmissionTime   time.Time `json:"submission_time"`
}
