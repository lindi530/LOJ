package constants

import "errors"

var (
	LoginUserIDKey  = "login_user_id"
	NoUserID        = int64(-1)
	NotLogin        = errors.New("not login")
	AccessUserIDKey = "access_user_id"
)
