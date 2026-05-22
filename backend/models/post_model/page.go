package post_model

import "GO1/models/cursor"

type PageInfo struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"` // form 使用 ShouldBindQuery
}

type CursorPostResp struct {
	Posts  []PostResponse
	Cursor cursor.CursorReq
}
