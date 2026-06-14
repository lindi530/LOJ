package course_model

type GetChapterInfoReq struct {
	CourseID  int64 `uri:"course_id" binding:"required"`
	ChapterID int64 `uri:"chapter_id" binding:"required"`
}

type GetChapterInfoResp struct {
	Problems []ChapterProblemInfo `json:"problems"`
	Video    ChapterVideoInfo     `json:"video"`
}

type ChapterProblemInfo struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type ChapterVideoInfo struct {
	ID int64 `json:"id"`
}
