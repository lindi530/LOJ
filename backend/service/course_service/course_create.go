package course_service

import (
	"GO1/database/mysql/course_mysql"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/course_model"
	"GO1/pkg/constants"
	pkg_image "GO1/pkg/image"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func CourseCreate(c *gin.Context, req *course_model.CourseCreateReq, cover *multipart.FileHeader) (resp response.Response) {
	userID, exists := c.Get(constants.LoginUserIDKey)
	if !exists {
		resp.Code = 1
		resp.Message = "用户ID错误"
		return
	}

	if req.Price < 0 {
		resp.Code = 1
		resp.Message = "课程价格错误"
		return
	}

	isFree := int8(1)
	if req.IsFree != nil && !*req.IsFree {
		isFree = 0
	}

	price := req.Price
	if isFree == 1 {
		price = 0
	}

	coverPath, err := saveCourseCover(c, cover)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		return
	}

	createdBy := userID.(int64)
	course := &course_model.Course{
		Title:       req.Title,
		Description: req.Description,
		CoverURL:    &coverPath,
		Price:       price,
		IsFree:      isFree,
		Status:      req.Status,
		Sort:        req.Sort,
		CreatedBy:   &createdBy,
	}

	if err := course_mysql.CreateCourse(course); err != nil {
		resp.Code = 1
		resp.Message = "课程创建失败"
		return
	}

	resp.Data = &course_model.CourseCreateResp{
		ID: course.ID,
	}
	return
}

func saveCourseCover(c *gin.Context, cover *multipart.FileHeader) (string, error) {
	if cover == nil {
		return "", errors.New("course cover is required")
	}

	if !inWriteList(filepath.Ext(cover.Filename)) {
		return "", errors.New("cover image format error")
	}

	imageLimitSize := float64(global.Config.Upload.Image.Size)
	if float64(cover.Size)/1024/1024 >= imageLimitSize {
		return "", errors.New("cover image size exceeds limit")
	}

	dirPath := global.Config.Upload.Image.Path
	if err := ensureDir(dirPath); err != nil {
		return "", err
	}

	uploadedImage, err := pkg_image.SaveUploadedImageByMD5Name(c, cover, dirPath)
	if err != nil {
		return "", errors.New("cover image save failed")
	}

	return uploadedImage.Path, nil
}

func ensureDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return err
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
