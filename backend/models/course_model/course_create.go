package course_model

type CourseCreateReq struct {
	Title       string  `form:"title" json:"title" binding:"required"`
	Description *string `form:"description" json:"description"`
	Price       float64 `form:"price" json:"price"`
	IsFree      *bool   `form:"isFree" json:"isFree"`
	Status      int8    `form:"status" json:"status"`
	Sort        int     `form:"sort" json:"sort"`
	// cover 图片
}

type CourseCreateResp struct {
	ID int64 `json:"id"`
}
