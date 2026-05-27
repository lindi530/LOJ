package competition_model

type CompetitionPlayer struct {
	ID            int64  `json:"id"`
	CompetitionID int64  `json:"competition_id"`
	UserID        int64  `json:"user_id"`
	UserName      string `json:"user_name"`
	IsCancel      bool   `json:"is_cancel"`
}
