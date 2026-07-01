package order_model

type PayCourseOrderURIReq struct {
	OrderNo string `uri:"order_no" binding:"required"`
}

type PayCourseOrderReq struct {
	PayType string `json:"payType" binding:"required"`
}

type PayCourseOrderResp struct {
	OrderNo  string  `json:"orderNo"`
	CourseID int64   `json:"courseId"`
	Status   int8    `json:"status"`
	Balance  float64 `json:"balance"`
	Owned    bool    `json:"owned"`
}
