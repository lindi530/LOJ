package user_post_service

import (
	"GO1/database/mysql/post_mysql"
	"GO1/middlewares/response"
	"GO1/models/cursor"
	"GO1/models/post_model"
	"GO1/pkg/constants"

	"github.com/gin-gonic/gin"
)

func GetNextCursorPosts(c *gin.Context, cursor cursor.CursorReq) (resp response.Response) {
	respData := &post_model.CursorPostResp{}
	if !cursor.HasMore {
		respData.Cursor.HasMore = false
		resp.Data = respData
		resp.Code = 1
		return
	}

	limit := 10
	posts, err := post_mysql.GetNextCursorPosts(cursor.Cursor, limit)
	if err != nil {
		resp.Code = 1
		resp.Message = "数据库查询失败"
		return
	}

	hasMore := len(posts) > limit
	respData.Cursor.HasMore = hasMore
	if hasMore {
		posts = posts[:limit]
	}
	respData.Cursor.Cursor = posts[len(posts)-1].PostID
	value, exists := c.Get(constants.LoginUserIDKey)
	if !exists {
		value = constants.NoUserID
	}
	respData.Posts = BuildPostsResponse(value.(int64), posts)
	resp.Data = respData
	return
}
