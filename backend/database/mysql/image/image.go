package image

import (
	"GO1/global"
	"GO1/models/upload_model"
	"errors"
	"time"

	"gorm.io/gorm"
)

func CheckImage(imgMD5Str string) bool {
	_, exists, err := GetImageByMD5(imgMD5Str)
	return err == nil && exists
}

func GetImageByMD5(imgMD5Str string) (upload_model.Image, bool, error) {
	var img upload_model.Image
	err := global.DB.Where("md5 = ?", imgMD5Str).First(&img).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return img, false, nil
	}
	if err != nil {
		return img, false, err
	}
	return img, true, nil
}

func CreateImage(imgMD5Str string, fileName string, filePath string) error {
	now := time.Now()
	img := upload_model.Image{
		MD5:       imgMD5Str,
		Name:      fileName,
		Path:      filePath,
		CreatedAt: now,
	}
	err := global.DB.Create(&img)
	return err.Error
}
