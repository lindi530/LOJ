package course_model

type CoursePayCallbackReq struct {
	OrderNo       string `json:"orderNo" binding:"required"`
	PayChannel    string `json:"payChannel"`
	TransactionID string `json:"transactionId" binding:"required"`
	PaidAt        string `json:"paidAt"`
}

type CoursePayCallbackResp struct {
	OrderNo string `json:"orderNo"`
	Status  int8   `json:"status"`
}
