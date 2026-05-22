package post_mysql

import (
	"GO1/global"
	"GO1/models/post_model"
)

func GetNextCursorPosts(cursor int64, limit int) ([]post_model.Post, error) {

	post := make([]post_model.Post, limit+1)

	if cursor == 0 {
		err := global.DB.Model(&post_model.Post{}).
			Preload("Author").
			Order("id DESC").
			Limit(limit + 1).
			Find(&post).Error
		return post, err
	}

	err := global.DB.Model(&post_model.Post{}).
		Where("id < ?", cursor).
		Preload("Author").
		Order("id DESC").
		Limit(limit + 1).
		Find(&post).Error
	return post, err
}
