package course_service

import (
	"GO1/database/mysql/course_mysql"
	"GO1/middlewares/response"
	"GO1/models/course_model"
)

func CourseList(_ *course_model.CourseListReq) (resp response.Response) {
	courses, total, err := course_mysql.GetCourseList()
	if err != nil {
		resp.Code = 1
		resp.Message = "数据查询错误"
		return
	}

	respData := &course_model.CourseListResp{
		List:  make([]course_model.CourseListRespItem, 0, len(courses)),
		Total: total,
	}

	for _, course := range courses {
		respData.List = append(respData.List, course_model.CourseListRespItem{
			ID:           course.ID,
			Title:        course.Title,
			Description:  course.Description,
			CoverURL:     course.CoverURL,
			Price:        course.Price,
			IsFree:       course.IsFree == 1,
			LessonCount:  course.LessonCount,
			StudentCount: course.StudentCount,
			CreatedBy:    course.CreatedBy,
		})
	}

	resp.Data = respData
	return
}
