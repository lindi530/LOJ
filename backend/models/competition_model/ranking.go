package competition_model

type GetRankingReq struct{}

type GetRankingResp struct {
	UserID               int64  `json:"user_id"`
	UserName             string `json:"user_name"`
	UserAvatar           string `json:"user_avatar"`
	UserRating           int    `json:"user_rating"` // 对应 “rank_score”
	UserCompetitionCount int    `json:"user_competition_count"`
}
