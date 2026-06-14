package upload_service

import (
	"os"
	"path"
	"path/filepath"
	"strconv"

	"GO1/database/mysql/course_mysql"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/upload_model"
)

func CreateVideoUploadTask() (resp response.Response) {
	asset, err := course_mysql.CreateEmptyVideoAsset()
	if err != nil {
		resp.Code = 1
		resp.Message = "视频上传任务创建失败"
		return
	}

	relativePath := path.Join(global.Config.Upload.Video.Path, strconv.FormatInt(asset.ID, 10))
	if err := os.MkdirAll(filepath.FromSlash(relativePath), 0755); err != nil {
		resp.Code = 1
		resp.Message = "视频上传目录创建失败"
		return
	}

	if err := course_mysql.UpdateVideoAssetURL(asset.ID, relativePath); err != nil {
		resp.Code = 1
		resp.Message = "视频上传任务更新失败"
		return
	}

	resp.Data = &upload_model.CreateVideoUploadTaskResp{
		UploadID:  asset.ID,
		ChunkSize: int(global.Config.Upload.Video.ChunkSize),
	}

	return
}
