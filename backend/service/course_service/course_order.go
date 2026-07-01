package course_service

import (
	"fmt"

	"GO1/database/mysql/course_mysql"
	"GO1/database/mysql/user_mysql"
	"GO1/middlewares/response"
	"GO1/models/course_model"
	pkg_snowflake "GO1/pkg/snowflake"
)

func CreateCourseOrder(userID int64, req *course_model.CreateCourseOrderReq) (resp response.Response) {
	if userID <= 0 {
		resp.Code = 1
		resp.Message = "need login"
		return
	}

	course, exist, err := course_mysql.GetCourseByID(req.CourseID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query course failed"
		return
	}
	if !exist {
		resp.Code = 1
		resp.Message = "course not found"
		return
	}

	balance, err := user_mysql.GetUserWalletBalance(userID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query user wallet failed"
		return
	}

	if course_mysql.IsCourseCreator(course, userID) {
		resp.Message = "course creator already has access"
		resp.Data = &course_model.CreateCourseOrderResp{
			CourseID: req.CourseID,
			Amount:   0,
			Balance:  balance,
			Status:   course_model.CourseOrderStatusPaid,
			Owned:    true,
			Free:     course.IsFree == 1 || course.Price <= 0,
		}
		return
	}

	owned, err := course_mysql.HasActiveCourseEnrollment(userID, req.CourseID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query course enrollment failed"
		return
	}
	if owned {
		resp.Message = "已拥有"
		resp.Data = &course_model.CreateCourseOrderResp{
			CourseID: req.CourseID,
			Amount:   course.Price,
			Balance:  balance,
			Status:   course_model.CourseOrderStatusPaid,
			Owned:    true,
			Free:     course.IsFree == 1 || course.Price <= 0,
		}
		return
	}

	if course.IsFree == 1 || course.Price <= 0 {
		if _, err := course_mysql.EnrollFreeCourse(userID, req.CourseID); err != nil {
			resp.Code = 1
			resp.Message = "create course enrollment failed"
			return
		}
		resp.Message = "已拥有"
		resp.Data = &course_model.CreateCourseOrderResp{
			CourseID: req.CourseID,
			Amount:   0,
			Balance:  balance,
			Status:   course_model.CourseOrderStatusPaid,
			Owned:    true,
			Free:     true,
		}
		return
	}

	orderNo := fmt.Sprintf("COURSE%d", pkg_snowflake.Snowflake{}.GenID())
	order, reused, err := course_mysql.GetOrCreatePendingCourseOrder(userID, req.CourseID, course.Price, orderNo)
	if err != nil {
		resp.Code = 1
		resp.Message = "create course order failed"
		return
	}
	if reused {
		resp.Message = "订单已存在"
	} else {
		resp.Message = "订单创建成功"
	}
	resp.Data = &course_model.CreateCourseOrderResp{
		CourseID: req.CourseID,
		OrderNo:  order.OrderNo,
		Amount:   order.Amount,
		Balance:  balance,
		Status:   order.Status,
		Owned:    false,
		Free:     false,
	}
	return
}
