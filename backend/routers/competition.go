package routers

import (
	"GO1/api"
	"GO1/middlewares/jwt"

	"github.com/gin-gonic/gin"
)

func CompetitionRouters(router *gin.RouterGroup) {
	auth := router.Group("/competition")

	auth.GET("", api.ApiGroups.CompetitionAPI.GetCompetitions)
	auth.GET("/rank")

	auth.Use(jwt.JWTAuthMiddleware())
	{
		auth.POST("create", api.ApiGroups.CompetitionAPI.CreateCompetition)
	}
}
