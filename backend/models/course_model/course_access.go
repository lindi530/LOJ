package course_model

type GetCourseAccessReq struct {
	CourseID int64 `uri:"course_id" binding:"required"`
}

type GetCourseAccessResp struct {
	CourseID  int64 `json:"courseId"`
	Owned     bool  `json:"owned"`
	Creator   bool  `json:"creator"`
	CanAccess bool  `json:"canAccess"`
}
