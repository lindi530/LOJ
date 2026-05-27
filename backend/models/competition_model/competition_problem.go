package competition_model

type CompetitionProblem struct {
	ID            int64 `json:"id"`
	CompetitionID int64 `json:"competition_id"`
	ProblemID     int   `json:"problem_id"`
}
