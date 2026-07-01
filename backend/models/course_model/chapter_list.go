package course_model

type GetChapterListReq struct {
	CourseID int64 `uri:"course_id" binding:"required"`
}

type GetChapterListResp struct {
	List      []CourseChapter `json:"list"`
	HasAccess bool            `json:"hasAccess"`
}
