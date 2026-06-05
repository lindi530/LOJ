package competition_model

type CompetitionSummary struct {
	ID               int64 `json:"id"`
	UserID           int64 `json:"user_id"`
	RankScore        int   `json:"rank_score"` // 对应user_rating
	CompetitionCount int   `json:"competition_count"`
}
