package routers

import (
	"GO1/api"
	"GO1/middlewares/jwt"

	"github.com/gin-gonic/gin"
)

func CalendarRouter(router *gin.RouterGroup) {
	auth := router.Group("/calendar")

	auth.POST("/submissions/:userId", api.ApiGroups.CalendarAPI.GetSubmissionCalendar)

	auth.Use(jwt.JWTAuthMiddleware())
	{

	}

}
