package order_service

import (
	"errors"
	"strings"

	"GO1/database/mysql/order_mysql"
	"GO1/middlewares/response"
	"GO1/models/order_model"
)

func PayCourseOrderByCoin(userID int64, orderNo string, req *order_model.PayCourseOrderReq) (resp response.Response) {
	if userID <= 0 {
		resp.Code = 1
		resp.Message = "need login"
		return
	}
	if strings.TrimSpace(orderNo) == "" {
		resp.Code = 1
		resp.Message = "order no is required"
		return
	}
	if req == nil || strings.TrimSpace(req.PayType) != "coin" {
		resp.Code = 1
		resp.Message = "支付类型错误"
		return
	}

	order, balance, owned, err := order_mysql.PayCourseOrderByCoin(userID, orderNo)
	if err != nil {
		resp.Code = 1
		switch {
		case errors.Is(err, order_mysql.ErrCourseOrderNotFound):
			resp.Message = "订单不存在"
		case errors.Is(err, order_mysql.ErrCourseOrderUserInvalid):
			resp.Message = "订单不属于当前用户"
		case errors.Is(err, order_mysql.ErrCourseOrderStatusInvalid):
			resp.Message = "订单状态异常"
		case errors.Is(err, order_mysql.ErrUserWalletNotFound):
			resp.Message = "虚拟币余额不足"
		case errors.Is(err, order_mysql.ErrCoinBalanceInsufficient):
			resp.Message = "虚拟币余额不足"
		default:
			resp.Message = "支付失败"
		}
		return
	}

	resp.Message = "支付成功"
	resp.Data = &order_model.PayCourseOrderResp{
		OrderNo:  order.OrderNo,
		CourseID: order.CourseID,
		Status:   order.Status,
		Balance:  balance,
		Owned:    owned,
	}
	return
}
