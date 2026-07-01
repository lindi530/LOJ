package course_service

import (
	"GO1/database/mysql/course_mysql"
	"GO1/middlewares/response"
	"GO1/models/course_model"
)

func ChapterList(userID int64, req *course_model.GetChapterListReq) (resp response.Response) {
	chapters, err := course_mysql.GetChapterList(req.CourseID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query chapter list failed"
		return
	}

	hasAccess := false
	if userID > 0 {
		hasAccess, err = course_mysql.UserCanAccessCourse(userID, req.CourseID)
		if err != nil {
			resp.Code = 1
			resp.Message = "query course access failed"
			return
		}
	}

	resp.Data = &course_model.GetChapterListResp{
		List:      chapters,
		HasAccess: hasAccess,
	}
	return
}
