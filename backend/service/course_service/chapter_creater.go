package course_service

import (
	"GO1/database/mysql/course_mysql"
	"GO1/middlewares/response"
	"GO1/models/course_model"
	"GO1/pkg/constants"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ChapterCreate(c *gin.Context, req *course_model.ChapterCreateReq) (resp response.Response) {
	userID, exists := c.Get(constants.LoginUserIDKey)
	loginUserID, ok := userID.(int64)
	if !exists || !ok || loginUserID == constants.NoUserID {
		resp.Code = 1
		resp.Message = "用户ID错误"
		return
	}

	title := strings.TrimSpace(req.Title)
	if req.CourseID <= 0 {
		resp.Code = 1
		resp.Message = "课程ID错误"
		return
	}
	if title == "" {
		resp.Code = 1
		resp.Message = "章节标题不能为空"
		return
	}
	if req.Sort <= 0 {
		resp.Code = 1
		resp.Message = "章节排序错误"
		return
	}
	if len(req.Problems) == 0 {
		resp.Code = 1
		resp.Message = "章节题目不能为空"
		return
	}
	if req.Video.ID <= 0 || strings.TrimSpace(req.Video.URL) == "" {
		resp.Code = 1
		resp.Message = "视频信息错误"
		return
	}

	now := time.Now()
	chapter := &course_model.CourseChapter{
		CourseID:    req.CourseID,
		Title:       title,
		Description: req.Description,
		Sort:        req.Sort,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	problems := make([]course_model.ChapterProblem, 0, len(req.Problems))
	for _, problem := range req.Problems {
		if problem.ID <= 0 {
			resp.Code = 1
			resp.Message = "题目ID错误"
			return
		}
		if problem.Sort <= 0 {
			resp.Code = 1
			resp.Message = "题目排序错误"
			return
		}

		problems = append(problems, course_model.ChapterProblem{
			CourseID:  req.CourseID,
			ProblemID: problem.ID,
			Sort:      problem.Sort,
		})
	}

	video := &course_model.ChapterVideo{
		CourseID: req.CourseID,
		VideoID:  req.Video.ID,
	}

	if err := course_mysql.CreateChapter(chapter, problems, video); err != nil {
		resp.Code = 1
		resp.Message = "章节创建失败"
		return
	}

	resp.Data = &course_model.ChapterCreateResp{}
	return
}
