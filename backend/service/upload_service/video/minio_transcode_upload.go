package video

import (
	"context"
	"errors"
	"io/fs"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	"GO1/global"

	miniogo "github.com/minio/minio-go/v7"
)

func uploadHLSFiles(ctx context.Context, bucket string, localDir string, objectPrefix string) error {
	return filepath.WalkDir(localDir, func(localPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(localDir, localPath)
		if err != nil {
			return err
		}
		objectName := path.Join(objectPrefix, filepath.ToSlash(relativePath))

		return uploadLocalFileToMinio(ctx, bucket, objectName, localPath, transcodedContentType(objectName))
	})
}

func uploadLocalFileToMinio(ctx context.Context, bucket string, objectName string, filePath string, contentType string) error {
	if global.MinioCore == nil {
		return errors.New("minio core unavailable")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if contentType == "" {
		contentType = transcodedContentType(objectName)
	}

	_, err = global.MinioCore.PutObject(
		ctx,
		bucket,
		objectName,
		file,
		fileInfo.Size(),
		"",
		"",
		miniogo.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func transcodedContentType(objectName string) string {
	switch strings.ToLower(path.Ext(objectName)) {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/mp2t"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	}

	contentType := mime.TypeByExtension(strings.ToLower(path.Ext(objectName)))
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}
