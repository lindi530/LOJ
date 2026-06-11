package upload_service

import (
	"GO1/global"
	"GO1/models/upload_model"
	pkg_image "GO1/pkg/image"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context, file []*multipart.FileHeader, rsp *[]upload_model.ResponseUploadImages) {
	dirPath := global.Config.Upload.Image.Path
	imageLimitSize := float64(global.Config.Upload.Image.Size)
	ensureDir(dirPath)

	for _, file := range file {
		filename := file.Filename
		if !inWriteList(filepath.Ext(filename)) {
			addResponse(rsp, filename, "上传失败", "文件格式错误")
			continue
		}

		size := float64(file.Size) / 1024 / 1024
		if size >= imageLimitSize {
			addResponse(rsp, filename, "上传失败", fmt.Sprintf("图片大小为：%.2fM，不能超过：%.2fM", size, imageLimitSize))
			continue
		}

		filePath := dirPath + "/" + filename
		uploadedImage, err := pkg_image.SaveUploadedImageByMD5(c, file, filePath)
		if err != nil {
			global.Logger.Error(err)
			addResponse(rsp, filename, "上传失败", "保存失败")
			continue
		}

		if uploadedImage.Existed {
			addResponse(rsp, filename, "上传成功", "图片已存在", uploadedImage.Path)
			continue
		}
		addResponse(rsp, filename, "上传成功", "保存成功", uploadedImage.Path)
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

func addResponse(rsp *[]upload_model.ResponseUploadImages, filename string, result string, msg string, paths ...string) {
	path := ""
	if len(paths) > 0 {
		path = paths[0]
	}

	*rsp = append(*rsp, upload_model.ResponseUploadImages{
		FileName: filename,
		Result:   result,
		Msg:      msg,
		Path:     path,
	})
}
