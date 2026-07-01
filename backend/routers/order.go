package routers

import (
	"GO1/api"
	"GO1/middlewares/jwt"

	"github.com/gin-gonic/gin"
)

func OrderRouter(router *gin.RouterGroup) {
	order := router.Group("/order")
	order.Use(jwt.JWTAuthMiddleware())
	{
		order.POST("/:order_no/pay", api.ApiGroups.OrderAPI.PayCourseOrder)
	}
}
