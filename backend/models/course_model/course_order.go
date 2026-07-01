package course_model

type CreateCourseOrderReq struct {
	CourseID int64 `uri:"course_id" binding:"required"`
}

type CreateCourseOrderResp struct {
	CourseID int64   `json:"courseId"`
	OrderNo  string  `json:"orderNo"`
	Amount   float64 `json:"amount"`
	Balance  float64 `json:"balance"`
	Status   int8    `json:"status"`
	Owned    bool    `json:"owned"`
	Free     bool    `json:"free"`
}
