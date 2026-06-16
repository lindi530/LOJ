package video

import (
	"context"
	"path"
	"strconv"

	"GO1/database/mysql/course_mysql"
	"GO1/database/redis/course_redis"
	"GO1/global"
	"GO1/middlewares/response"
	"GO1/models/upload_model"

	miniogo "github.com/minio/minio-go/v7"
)

func CreateVideoUploadTask(ctx context.Context) (resp response.Response) {
	asset, err := course_mysql.CreateEmptyVideoAsset()
	if err != nil {
		resp.Code = 1
		resp.Message = "视频上传任务创建失败"
		return
	}

	bucket := global.Config.Minio.Bucket.Video
	prefix := videoSourceObjectPrefix()
	tempObject := path.Join(prefix, strconv.FormatInt(asset.ID, 10)+".upload")
	minioUploadID, err := global.MinioCore.NewMultipartUpload(ctx, bucket, tempObject, miniogo.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		resp.Code = 1
		resp.Message = "视频上传任务创建失败"
		return
	}

	if err := course_redis.SaveVideoMultipartUpload(ctx, asset.ID, course_redis.VideoMultipartUploadMeta{
		MinioUploadID: minioUploadID,
		Bucket:        bucket,
		TempObject:    tempObject,
	}); err != nil {
		_ = global.MinioCore.AbortMultipartUpload(context.Background(), bucket, tempObject, minioUploadID)
		resp.Code = 1
		resp.Message = "视频上传任务保存失败"
		return
	}

	originPath := path.Join(global.Config.Upload.Video.Path, strconv.FormatInt(asset.ID, 10))
	if err := course_mysql.UpdateVideoAssetOriginPath(asset.ID, originPath); err != nil {
		_ = global.MinioCore.AbortMultipartUpload(context.Background(), bucket, tempObject, minioUploadID)
		resp.Code = 1
		resp.Message = "视频上传任务更新失败"
		return
	}

	resp.Data = &upload_model.CreateVideoUploadTaskResp{
		UploadID:      asset.ID,
		MinioUploadID: minioUploadID,
		ChunkSize:     int(global.Config.Upload.Video.ChunkSize),
	}

	return
}
