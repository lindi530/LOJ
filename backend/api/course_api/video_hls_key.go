package course_api

import (
	"GO1/global"
	"net/http"

	"GO1/middlewares/response"
	"GO1/models/course_model"
	"GO1/service/context"
	"GO1/service/course_service"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) GetVideoHLSKey(c *gin.Context) {
	var req course_model.GetVideoHLSKeyReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	userID, _ := context.GetContext(c, context.CtxUserIdKey).(int64)
	key, status, message := course_service.GetVideoHLSKey(c.Request.Context(), userID, &req)

	global.Logger.Info("message:", message)
	if status != http.StatusOK {
		c.AbortWithStatusJSON(status, gin.H{
			"code":    1,
			"message": message,
		})
		return
	}

	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/octet-stream", key)
}
