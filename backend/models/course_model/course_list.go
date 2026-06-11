package course_model

type CourseListReq struct {
}

type CourseListResp struct {
	List  []CourseListRespItem `json:"list"`
	Total int64                `json:"total"`
}

type CourseListRespItem struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	CoverURL     *string `json:"coverUrl"`
	Price        float64 `json:"price"`
	IsFree       bool    `json:"isFree"`
	LessonCount  int     `json:"lessonCount"`
	StudentCount int     `json:"studentCount"`
}
