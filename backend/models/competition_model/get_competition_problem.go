package competition_model

type GetCompetitionProblemsReq struct {
	CompetitionID int64 `json:"competition_id" uri:"competition_id" binding:"required"`
}

type GetCompetitionProblemsResp struct {
	ProblemNumber string `json:"problem_number"`
	ProblemID     int    `json:"problem_id"`
	ProblemTitle  string `json:"problem_title"`
	AcCount       int    `json:"ac_count"`
	SubmitCount   int    `json:"submit_count"`
	HasAccess     bool   `json:"has_access"`
}
