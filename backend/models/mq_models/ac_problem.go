package mq_models

import "time"

type ACProblem struct {
	UserID  int64     `json:"user_id"`
	AcCount int       `json:"ac_count"`
	Date    time.Time `json:"date"`
}
