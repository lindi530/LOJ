package course_service

import (
	"errors"

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
			Resolution:   profile.Resolution,
			PlaylistPath: profile.PlaylistPath,
			FileSize:     profile.FileSize,
		}
	}

	resp.Data = &course_model.GetChapterVideoResp{
		ID:        videoAsset.ID,
		CoverPath: videoAsset.CoverPath,
		Profiles:  respProfiles,
	}
	return
}
