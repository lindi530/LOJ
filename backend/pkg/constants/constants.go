package constants

import "errors"

var (
	LoginUserIDKey   = "login_user_id"
	LoginUserNameKey = "login_user_name"
	NoUserID         = int64(-1)
	NoUserName       = ""
	NotLogin         = errors.New("not login")
	AccessUserIDKey  = "access_user_id"
)
