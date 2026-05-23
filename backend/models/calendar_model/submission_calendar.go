package calendar_model

import (
	"time"
)

type CalendarSubmission struct {
	ID      int64     `json:"id"`
	UserId  int64     `json:"user_id"`
	ACCount int       `json:"ac_count"`
	ACDate  time.Time `json:"ac_date"`
}

func (CalendarSubmission) TableName() string {
	return "calendar_submissions"
}

type CalendarSubmissionReq struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type CalendarSubmissionResp struct {
	TotalAC    int              `json:"total_ac"`
	ActiveDays int              `json:"active_days"`
	Days       []SubmissionDays `json:"days"`
}

type SubmissionDays struct {
	Date    time.Time `json:"data"`
	ACCount int       `json:"ac_count"`
}
