package course_service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"time"

	"GO1/database/mysql/course_mysql"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/course_model"

	"gorm.io/gorm"
)

func CoursePayCallback(rawBody []byte, signature string, req *course_model.CoursePayCallbackReq) (resp response.Response) {
	if !verifyCoursePayCallbackSign(rawBody, signature) {
		resp.Code = 1
		resp.Message = "invalid pay sign"
		return
	}

	paidAt, err := parseCoursePaidAt(req.PaidAt)
	if err != nil {
		resp.Code = 1
		resp.Message = "invalid paidAt"
		return
	}

	order, err := course_mysql.PayCourseOrder(req.OrderNo, req.PayChannel, req.TransactionID, paidAt)
	if err != nil {
		resp.Code = 1
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			resp.Message = "order not found"
		case errors.Is(err, course_mysql.ErrCourseOrderStatusInvalid):
			resp.Message = "order status invalid"
		default:
			resp.Message = "pay callback failed"
		}
		return
	}

	resp.Message = "success"
	resp.Data = &course_model.CoursePayCallbackResp{
		OrderNo: order.OrderNo,
		Status:  course_model.CourseOrderStatusPaid,
	}
	return
}

func verifyCoursePayCallbackSign(rawBody []byte, signature string) bool {
	secret := coursePayCallbackSecret()
	signature = strings.TrimSpace(signature)
	if secret == "" || signature == "" {
		return false
	}

	signature = strings.TrimPrefix(strings.ToLower(signature), "sha256=")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(rawBody)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func coursePayCallbackSecret() string {
	if secret := strings.TrimSpace(os.Getenv("COURSE_PAY_CALLBACK_SECRET")); secret != "" {
		return secret
	}
	if global.Config == nil {
		return ""
	}
	return strings.TrimSpace(global.Config.CoursePay.CallbackSecret)
}

func parseCoursePaidAt(input string) (time.Time, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return time.Now(), nil
	}

	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05"} {
		if paidAt, err := time.Parse(layout, input); err == nil {
			return paidAt, nil
		}
	}

	return time.Time{}, errors.New("invalid paidAt")
}
