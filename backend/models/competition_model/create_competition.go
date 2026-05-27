package competition_model

import "time"

type CompetitionReq struct {
	Name         string    `json:"name"`
	Visibility   bool      `json:"visibility"`
	Password     string    `json:"password"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Announcement string    `json:"announcement"`
	ProblemIDs   []int     `json:"problem_ids"`
}
