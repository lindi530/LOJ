package competition_model

type RankListResp struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Rating   int    `json:"ranking"`
	Avatar   string `json:"avatar"`
}
