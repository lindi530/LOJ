package image

import (
	mysql_image "GO1/database/mysql/image"
	"GO1/pkg/md5"
	"io"
	"mime/multipart"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var WriteImageList = []string{
	".jpg",
	".png",
	".gif",
	".jpeg",
	".webp",
}

type UploadedImage struct {
	MD5     string
	Path    string
	Existed bool
}

func SaveUploadedImageByMD5(c *gin.Context, file *multipart.FileHeader, filePath string) (*UploadedImage, error) {
	return saveUploadedImageByMD5(c, file, func(fileMD5 string) string {
		return filePath
	})
}

func SaveUploadedImageByMD5Name(c *gin.Context, file *multipart.FileHeader, dirPath string) (*UploadedImage, error) {
	return saveUploadedImageByMD5(c, file, func(fileMD5 string) string {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		return path.Join(dirPath, fileMD5+ext)
	})
}

func saveUploadedImageByMD5(c *gin.Context, file *multipart.FileHeader, buildPath func(string) string) (*UploadedImage, error) {
	fileObj, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileObj.Close()

	byteData, err := io.ReadAll(fileObj)
	if err != nil {
		return nil, err
	}

	fileMD5 := md5.MD5(byteData)
	existedImage, exists, err := mysql_image.GetImageByMD5(fileMD5)
	if err != nil {
		return nil, err
	}
	if exists {
		return &UploadedImage{
			MD5:     fileMD5,
			Path:    existedImage.Path,
			Existed: true,
		}, nil
	}

	filePath := buildPath(fileMD5)
	if err := mysql_image.CreateImage(fileMD5, file.Filename, filePath); err != nil {
		return nil, err
	}
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return nil, err
	}

	return &UploadedImage{
		MD5:  fileMD5,
		Path: filePath,
	}, nil
}
