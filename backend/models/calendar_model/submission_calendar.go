package calendar_model

import (
	"time"
)

type CalendarSubmission struct {
	ID      int64     `json:"id"`
	UserID  int64     `json:"user_id"`
	ACCount int       `json:"ac_count"`
	ACDate  time.Time `json:"ac_date"`
}

func (CalendarSubmission) TableName() string {
	return "calendar_submissions"
}

type CalendarSubmissionReq struct {
	StartDate time.Time `form:"start_date" json:"start_date" binding:"required"`
	EndDate   time.Time `form:"end_date" json:"end_date" binding:"required"`
}

type CalendarSubmissionResp struct {
	TotalAC    int              `json:"total_ac"`
	ActiveDays int              `json:"active_days"`
	Days       []SubmissionDays `json:"days"`
}

type SubmissionDays struct {
	Date    time.Time `json:"date"`
	ACCount int       `json:"ac_count"`
}
