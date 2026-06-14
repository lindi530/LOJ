package course_api

import (
	"GO1/middlewares/response"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) GetChapterVideo(c *gin.Context) {
	response.OkWithMessage("成功获取视频", c)
}
