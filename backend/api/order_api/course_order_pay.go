package order_api

import (
	"GO1/middlewares/response"
	"GO1/models/order_model"
	service_context "GO1/service/context"
	"GO1/service/order_service"

	"github.com/gin-gonic/gin"
)

func (OrderAPI) PayCourseOrder(c *gin.Context) {
	var uriReq order_model.PayCourseOrderURIReq
	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	var req order_model.PayCourseOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(response.BadRequest, c)
		return
	}

	userID, ok := service_context.GetContext(c, service_context.CtxUserIdKey).(int64)
	if !ok || userID <= 0 {
		response.FailWithCode(response.NeedLogin, c)
		return
	}

	resp := order_service.PayCourseOrderByCoin(userID, uriReq.OrderNo, &req)
	if resp.Code == 1 {
		response.FailWithMessage(resp.Message, c)
		return
	}

	response.Ok(resp.Data, resp.Message, c)
}
