package course_service

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"GO1/database/mysql/course_mysql"
	"GO1/middlewares/response"
	"GO1/models/course_model"

	"gorm.io/gorm"
)

func GetChapterVideo(req *course_model.GetChapterVideoReq) (resp response.Response) {
	if req.VideoID <= 0 {
		resp.Code = 1
		resp.Message = "invalid video_id"
		return
	}

	exists, err := course_mysql.ChapterVideoExists(req.CourseID, req.ChapterID, req.VideoID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query chapter video failed"
		return
	}
	if !exists {
		resp.Code = 1
		resp.Message = "video not found"
		return
	}

	videoAsset, err := course_mysql.GetVideoAssetByID(req.VideoID)
	if err != nil {
		resp.Code = 1
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Message = "video not found"
			return
		}
		resp.Message = "query video asset failed"
		return
	}

	profiles, err := course_mysql.GetVideoProfilesByVideoID(req.VideoID)
	if err != nil {
		resp.Code = 1
		resp.Message = "query video profiles failed"
		return
	}

	respProfiles := make([]course_model.GetChapterVideoProfile, len(profiles))
	for i, profile := range profiles {
		respProfiles[i] = course_model.GetChapterVideoProfile{
			Resolution: profile.Resolution,
			Width:      profile.Width,
			Height:     profile.Height,
			Bitrate:    profile.Bitrate,
			PlayURL:    courseVideoProfileURL(req.VideoID, profile.PlaylistPath),
			FileSize:   profile.FileSize,
		}
	}

	sizeBytes := int64(0)
	if videoAsset.SizeBytes != nil {
		sizeBytes = *videoAsset.SizeBytes
	}
	coverURL := ""
	if videoAsset.CoverPath != "" {
		coverURL = courseVideoCoverURL(req.VideoID)
	}

	resp.Data = &course_model.GetChapterVideoResp{
		ID:          videoAsset.ID,
		Title:       "",
		Description: "",
		CoverURL:    coverURL,
		Duration:    videoAsset.Duration,
		SizeBytes:   sizeBytes,
		Profiles:    respProfiles,
	}
	return
}

func courseVideoCoverURL(videoID int64) string {
	return fmt.Sprintf(
		"/video-covers/%d.jpg",
		videoID,
	)
}

func courseVideoProfileURL(videoID int64, playlistPath string) string {
	normalizedPath := strings.TrimSpace(strings.ReplaceAll(playlistPath, "\\", "/"))
	if normalizedPath == "" {
		return ""
	}

	if parsedURL, err := url.Parse(normalizedPath); err == nil && parsedURL.IsAbs() {
		normalizedPath = parsedURL.Path
	}

	normalizedPath = strings.TrimLeft(normalizedPath, "/")
	videoIDText := strconv.FormatInt(videoID, 10)

	relativePath := normalizedPath
	if markerIndex := strings.Index("/"+normalizedPath, "/"+videoIDText+"/"); markerIndex >= 0 {
		relativePath = ("/" + normalizedPath)[markerIndex+len(videoIDText)+2:]
	}
	relativePath = strings.TrimLeft(relativePath, "/")
	if relativePath == "" {
		return ""
	}

	return "/" + path.Join("videos", videoIDText, relativePath)
}
