package user_service

import (
	"GO1/database/mysql/user_mysql"
	"GO1/models/user_model"
)

func Register(register user_model.ParamRegister) (bool, error) {
	result := user_mysql.CheckUser(user_mysql.UserNameParam(register.Name))

	if result {
		return false, nil
	}
	if err := user_mysql.Register(register); err != nil {
		return false, err
	}
	return true, nil
}
