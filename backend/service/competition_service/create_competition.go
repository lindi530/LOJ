package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"
	"GO1/pkg/password_hash"

	"github.com/gin-gonic/gin"
)

func CreateCompetition(c *gin.Context, req *competition_model.CreateCompetitionReq) (resp response.Response) {
	userID, exists := c.Get(constants.LoginUserIDKey)
	if !exists {
		resp.Code = 1
		resp.Message = "用户ID错误"
		return
	}
	userName, exists := c.Get(constants.LoginUserNameKey)
	if !exists {
		resp.Code = 1
		resp.Message = "用户名不匹配"
		return
	}

	passwordHash, err := password_hash.HashPassword(req.Password)
	if err != nil {
		resp.Code = 1
		resp.Message = "密码加密错误"
		return
	}

	competition := &competition_model.Competition{
		Name:         req.Name,
		Status:       1,
		Visibility:   req.Visibility,
		PasswordHash: passwordHash,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		CreatorID:    userID.(int64),
		CreatorName:  userName.(string),
		Announcement: req.Announcement,
	}

	if err := competition_mysql.CreateCompetition(competition, req.Problems); err != nil {
		resp.Code = 1
		resp.Message = "竞赛创建失败"
		return
	}

	resp.Message = "创建成功"
	return
}
