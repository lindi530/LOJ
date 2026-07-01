package user_mysql

import (
	"GO1/global"
	"GO1/models/user_model"
	"GO1/pkg/password_hash"
	"GO1/pkg/snowflake"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func Register(register user_model.ParamRegister) error {

	userId := snowflake.Snowflake{}.GenID()
	hashPassword, err := password_hash.HashPassword(register.Password)

	if err != nil {
		global.Logger.Error(fmt.Sprintf("%s hash password error: %v", register.Name, err))
		return err
	}
	now := time.Now()
	user := user_model.User{
		UserID:         userId,
		UserName:       register.Name,
		Password:       hashPassword,
		Email:          register.Email,
		Gender:         global.Config.Settings.User.Gender,
		Quote:          global.Config.Settings.User.Quote,
		Avatar:         global.Config.Settings.User.Avatar,
		FollowerCount:  0,
		FollowingCount: 0,
		CreateTime:     now,
		UpdateTime:     now,
	}

	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return EnsureUserWalletWithTx(tx, userId)
	}); err != nil {
		global.Logger.Error(fmt.Sprintf("%s register error: %v", register.Name, err))
		return err
	}

	global.Logger.Info(user)
	return nil
}
