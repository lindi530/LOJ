package cursor

type CursorReq struct {
	Cursor  int64 `form:"cursor" json:"cursor"`
	HasMore bool  `form:"has_more" json:"has_more"`
}
