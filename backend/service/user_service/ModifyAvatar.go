package user_service

import (
	"GO1/database/mysql/user_mysql"
	"GO1/global"
	"GO1/models/user_model"
	pkg_image "GO1/pkg/image"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func ModifyAvatar(c *gin.Context, userId int64, avatar *multipart.FileHeader) gin.H {
	if !AuthUser(c, userId) {
		global.Logger.Error("无权限")
		return gin.H{
			"avatarPath": "",
			"msg":        "无权限",
		}
	}

	if !inWriteList(filepath.Ext(avatar.Filename)) {
		return gin.H{
			"avatarPath": "",
			"msg":        "格式不允许",
		}
	}

	imageLimitSize := float64(global.Config.Upload.Image.Size)
	if float64(avatar.Size/1024/1024) > imageLimitSize {
		return gin.H{
			"avatarPath": "",
			"msg":        "超出限制大小",
		}
	}

	dirPath := global.Config.Upload.Image.Path
	ensureDir(dirPath)

	filePath := dirPath + "/" + avatar.Filename
	uploadedImage, err := pkg_image.SaveUploadedImageByMD5(c, avatar, filePath)
	if err != nil {
		return gin.H{
			"avatarPath": "",
			"msg":        "文件保存失败",
		}
	}

	user_mysql.ModifyAvatar(userId, uploadedImage.MD5)
	avatarPath := user_model.GetAvatarPath(uploadedImage.MD5)
	return gin.H{
		"avatar": avatarPath,
		"msg":    "上传成功",
	}
}

func ensureDir(dirPath string) {
	_, err := os.Stat(dirPath)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		os.MkdirAll(dirPath, 0755)
	}
}

func inWriteList(ext string) bool {
	ext = strings.ToLower(ext)
	for _, EXT := range pkg_image.WriteImageList {
		if EXT == ext {
			return true
		}
	}
	return false
}
