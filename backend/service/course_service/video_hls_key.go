package course_service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"GO1/database/mysql/course_mysql"
	"GO1/global"
	"GO1/models/course_model"

	miniogo "github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

const hlsKeyFileName = "enc.key"

func GetVideoHLSKey(ctx context.Context, userID int64, req *course_model.GetVideoHLSKeyReq) ([]byte, int, string) {
	if req.VideoID <= 0 || strings.TrimSpace(req.Variant) == "" {
		return nil, http.StatusBadRequest, "invalid hls key request"
	}

	videoAsset, err := course_mysql.GetVideoAssetByID(req.VideoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, "video not found"
		}
		return nil, http.StatusInternalServerError, "query video asset failed"
	}
	if videoAsset.Status != course_model.VideoAssetStatusPlayable {
		return nil, http.StatusNotFound, "video not playable"
	}

	profile, err := course_mysql.GetVideoProfileByResolution(req.VideoID, req.Variant)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, "video profile not found"
		}
		return nil, http.StatusInternalServerError, "query video profile failed"
	}

	canAccess, err := course_mysql.UserCanAccessVideoCourse(userID, req.VideoID)
	if err != nil {
		return nil, http.StatusInternalServerError, "query course access failed"
	}
	if !canAccess {
		return nil, http.StatusForbidden, "no course access"
	}

	variantDir := hlsVariantDirByResolution(req.Variant)
	if variantDir == "" {
		variantDir = hlsVariantDirFromPlaylist(req.VideoID, profile.PlaylistPath)
	}
	if variantDir == "" {
		return nil, http.StatusInternalServerError, "video profile path invalid"
	}

	key, err := readVideoHLSKey(ctx, req.VideoID, variantDir)
	if err != nil {
		return nil, http.StatusInternalServerError, "read hls key failed"
	}
	if len(key) != 16 {
		return nil, http.StatusInternalServerError, "invalid hls key"
	}

	return key, http.StatusOK, ""
}

func readVideoHLSKey(ctx context.Context, videoID int64, variantDir string) ([]byte, error) {
	if global.MinioCore == nil {
		return nil, fmt.Errorf("minio core unavailable")
	}

	objectName := path.Join(videoHLSKeyObjectPrefix(videoID), variantDir, hlsKeyFileName)
	reader, _, _, err := global.MinioCore.GetObject(
		ctx,
		global.Config.Minio.Bucket.Video,
		objectName,
		miniogo.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func hlsVariantDirByResolution(resolution string) string {
	resolution = strings.TrimSpace(resolution)
	for _, variant := range global.Config.Upload.Video.HLS.Variants {
		if strings.EqualFold(variant.Resolution, resolution) || strings.EqualFold(variant.Dir, resolution) {
			return strings.Trim(strings.ReplaceAll(variant.Dir, "\\", "/"), "/")
		}
	}
	return ""
}

func hlsVariantDirFromPlaylist(videoID int64, playlistPath string) string {
	normalizedPath := strings.TrimSpace(strings.ReplaceAll(playlistPath, "\\", "/"))
	if normalizedPath == "" {
		return ""
	}

	if parsedURL, err := url.Parse(normalizedPath); err == nil && parsedURL.IsAbs() {
		normalizedPath = parsedURL.Path
	}

	normalizedPath = strings.TrimLeft(normalizedPath, "/")
	videoIDText := strconv.FormatInt(videoID, 10)
	if markerIndex := strings.Index("/"+normalizedPath, "/"+videoIDText+"/"); markerIndex >= 0 {
		normalizedPath = ("/" + normalizedPath)[markerIndex+len(videoIDText)+2:]
	}

	variantDir := path.Dir(strings.TrimLeft(normalizedPath, "/"))
	if variantDir == "." || variantDir == "/" {
		return ""
	}
	return strings.Trim(variantDir, "/")
}

func videoHLSKeyObjectPrefix(videoID int64) string {
	return path.Join(videoHLSKeyRootPrefix(), strconv.FormatInt(videoID, 10))
}

func videoHLSKeyRootPrefix() string {
	prefix := cleanObjectPrefix(global.Config.Minio.Bucket.VideoPlayPath)
	if prefix == "" {
		return "hls-key"
	}

	base := path.Base(prefix)
	dir := path.Dir(prefix)
	keyBase := base + "-key"
	if dir == "." {
		return keyBase
	}
	return path.Join(dir, keyBase)
}

func cleanObjectPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	prefix = strings.ReplaceAll(prefix, "\\", "/")
	return strings.Trim(prefix, "/")
}
