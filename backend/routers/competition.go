package routers

import (
	"GO1/api"
	"GO1/middlewares/jwt"

	"github.com/gin-gonic/gin"
)

func CompetitionRouters(router *gin.RouterGroup) {
	auth := router.Group("/competition")

	auth.GET("/status", api.ApiGroups.CompetitionAPI.GetCompetitions)
	auth.GET("/:competition_id", api.ApiGroups.CompetitionAPI.GetCompetition)
	auth.GET("/rank")

	auth.Use(jwt.JWTAuthMiddleware())
	{
		auth.POST("/create", api.ApiGroups.CompetitionAPI.CreateCompetition)
		auth.POST("/:competition_id/enter", api.ApiGroups.CompetitionAPI.EnterCompetition)
		auth.GET("/:competition_id/has_entered", api.ApiGroups.CompetitionAPI.HasEnterCompetition)
		auth.GET("/:competition_id/problems", api.ApiGroups.CompetitionAPI.GetCompetitionProblems)
		auth.GET("/:competition_id/:problem_number", api.ApiGroups.CompetitionAPI.GetCompetitionProblemInfo)
		auth.POST("/:competition_id/submit/:problem_number", api.ApiGroups.CompetitionAPI.SubmitCompetitionProblem)
		auth.GET("/:competition_id/submissions", api.ApiGroups.CompetitionAPI.GetCompetitionSubmissions)
	}
}
