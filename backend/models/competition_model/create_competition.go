package competition_model

import "time"

type CreateCompetitionReq struct {
	Name         string                     `json:"name"`
	Visibility   bool                       `json:"visibility"`
	Password     string                     `json:"password"`
	StartTime    time.Time                  `json:"start_time"`
	EndTime      time.Time                  `json:"end_time"`
	Announcement string                     `json:"announcement"`
	Problems     []CreateCompetitionProblem `json:"problems"`
}

type CreateCompetitionProblem struct {
	ProblemNumber string `json:"problem_number"`
	ProblemID     int    `json:"problem_id"`
}
