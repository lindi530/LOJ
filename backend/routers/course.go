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
	course.Use(jwt.JWTAuthMiddleware())
	{
		course.POST("/create", api.ApiGroups.CourseAPI.CourseCreate)

	}

}
