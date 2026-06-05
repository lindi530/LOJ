package competition_model

type GetCompetitionHistoryReq struct {
	UserID int64 `json:"user_id"`
}

type GetCompetitionHistoryResp struct {
	RankScore       int    `json:"rank_score"`        // 当前rank分
	RankScoreChange int    `json:"rank_score_change"` // 本次比赛的变化值
	Rank            int    `json:"rank"`              // 本次比赛排名
	CompetitionId   int64  `json:"competition_id"`
	CompetitionName string `json:"competition_name"`
}
