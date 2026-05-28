package competition_model

type CompetitionProblem struct {
	ID            int64  `json:"id"`
	CompetitionID int64  `json:"competition_id"`
	ProblemNumber string `json:"problem_number"`
	ProblemID     int    `json:"problem_id"`
	AcCount       int    `json:"ac_count"`
	SubmitCount   int    `json:"submit_count"`
}
