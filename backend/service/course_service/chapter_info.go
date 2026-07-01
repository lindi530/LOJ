package course_service

import (
	"GO1/database/mysql/course_mysql"
	"GO1/middlewares/response"
	"GO1/models/course_model"
)

func ChapterInfo(userID int64, req *course_model.GetChapterInfoReq) (resp response.Response) {
	canAccess, err := course_mysql.UserCanAccessCourse(userID, req.CourseID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query course access failed"
		return
	}
	if !canAccess {
		resp.Code = 1
		resp.Message = "no course access"
		return
	}

	problems, err := course_mysql.GetChapterProblems(req.CourseID, req.ChapterID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query chapter problems failed"
		return
	}

	video, err := course_mysql.GetChapterVideo(req.CourseID, req.ChapterID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query chapter video failed"
		return
	}

	resp.Data = &course_model.GetChapterInfoResp{
		Problems: problems,
		Video:    video,
	}
	return
}
