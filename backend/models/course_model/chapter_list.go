package course_model

type GetChapterListReq struct {
	CourseID int64 `uri:"course_id" binding:"required"`
}
