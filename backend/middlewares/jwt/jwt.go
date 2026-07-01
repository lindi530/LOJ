package jwt

import (
	"errors"
	"strings"

	"GO1/middlewares/response"
	"GO1/pkg/jwt"
	"GO1/service/context"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.FailWithCode(response.NeedLogin, c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.FailWithCode(response.ExpiredAccessToken, c)
			} else {
				response.FailWithCode(response.InvalidAccessToken, c)
			}
			c.Abort()
			return
		}

		context.SetContext(c, context.CtxUserIdKey, claims.UserId)

		c.Next()
	}
}
