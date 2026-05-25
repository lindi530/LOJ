package context_script

import (
	"GO1/pkg/constants"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SaveAccessUserID(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	c.Set(constants.AccessUserIDKey, userID)
}
