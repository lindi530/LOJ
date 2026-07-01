package routers

import (
	"GO1/api"
	"GO1/middlewares/jwt"

	"github.com/gin-gonic/gin"
)

func CourseRouter(router *gin.RouterGroup) {
	course := router.Group("/course")

	course.GET("", api.ApiGroups.CourseAPI.CourseList)
	//course.GET("/:course_id", api.ApiGroups.CourseAPI.GetCourse)
	course.GET("/:course_id/chapter", jwt.JWTAuthMiddleware(), api.ApiGroups.CourseAPI.GetChapterList)
	course.POST("/order/pay-callback", api.ApiGroups.CourseAPI.CoursePayCallback)
	course.Use(jwt.JWTAuthMiddleware())
	{
		course.POST("/create", api.ApiGroups.CourseAPI.CourseCreate)
		course.POST("/:course_id/order", api.ApiGroups.CourseAPI.CreateCourseOrder)
		course.GET("/:course_id/access", api.ApiGroups.CourseAPI.GetCourseAccess)
		course.GET("/:course_id/chapter/:chapter_id", api.ApiGroups.CourseAPI.GetChapterInfo)
		course.GET("/:course_id/chapter/:chapter_id/:video_id", api.ApiGroups.CourseAPI.GetChapterVideo)
		course.GET("/video/:video_id/hls-key", api.ApiGroups.CourseAPI.GetVideoHLSKey)
	}

	chapter := course.Group("/chapter")
	chapter.Use(jwt.JWTAuthMiddleware())
	{
		//chapter.GET("", api.ApiGroups.CourseAPI.GetChapterList)

		chapter.POST("/create", api.ApiGroups.CourseAPI.ChapterCreate)
	}
}
