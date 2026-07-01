package course_api

import (
	"encoding/json"
	"io"
	"strings"

	"GO1/middlewares/response"
	"GO1/models/course_model"
	"GO1/service/course_service"

	"github.com/gin-gonic/gin"
)

func (CourseAPI) CoursePayCallback(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	var req course_model.CoursePayCallbackReq
	if err := json.Unmarshal(rawBody, &req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}
	if strings.TrimSpace(req.OrderNo) == "" || strings.TrimSpace(req.TransactionID) == "" {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	signature := c.GetHeader("X-Course-Pay-Sign")
	if signature == "" {
		signature = c.GetHeader("X-Pay-Sign")
	}

	resp := course_service.CoursePayCallback(rawBody, signature, &req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.Ok(resp.Data, resp.Message, c)
}
