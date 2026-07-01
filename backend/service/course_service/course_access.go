package course_service

import (
	"GO1/database/mysql/course_mysql"
	"GO1/middlewares/response"
	"GO1/models/course_model"
)

func CourseAccess(userID int64, req *course_model.GetCourseAccessReq) (resp response.Response) {
	if userID <= 0 {
		resp.Code = 1
		resp.Message = "need login"
		return
	}

	access, exist, err := course_mysql.GetUserCourseAccessStatus(userID, req.CourseID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query course access failed"
		return
	}
	if !exist {
		resp.Code = 1
		resp.Message = "course not found"
		return
	}

	resp.Data = &course_model.GetCourseAccessResp{
		CourseID:  access.CourseID,
		Owned:     access.Owned,
		Creator:   access.Creator,
		CanAccess: access.CanAccess,
	}
	return
}
