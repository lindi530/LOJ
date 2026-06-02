package competition_model

type CompetitionSubmission struct {
	ID uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`

	CompetitionID int64 `json:"competition_id" gorm:"column:competition_id"`
	ProblemID     int   `json:"problem_id" gorm:"column:problem_id"`
	SubmissionID  int64 `json:"submission_id" gorm:"column:submission_id"`
	UserID        int64 `json:"user_id" gorm:"column:user_id"`

	IsAC      bool `json:"is_ac" gorm:"column:is_ac"`
	IsFirstAC bool `json:"is_first_ac" gorm:"column:is_first_ac"`
}
