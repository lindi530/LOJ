package problem_submission_model

import "time"

type UserSubmission struct {
	Id          int64     `json:"id"`
	ProblemId   uint      `json:"problem_id"`
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Language    string    `json:"language"`
	State       string    `json:"state"`
	ExecTime    int       `json:"exec_time"`    // ms
	MemoryUsage int       `json:"memory_usage"` // KB
	Score       float64   `json:"score"`
	CreatedAt   time.Time `json:"created_at"`
}
