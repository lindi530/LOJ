package course_model

type GetChapterVideoReq struct {
	CourseID  int64 `uri:"course_id" binding:"required"`
	ChapterID int64 `uri:"chapter_id" binding:"required"`
	VideoID   int64 `uri:"video_id" binding:"required"`
}

type GetChapterVideoResp struct {
	ID        int64                    `json:"id"`
	CoverPath string                   `json:"cover_path"`
	Profiles  []GetChapterVideoProfile `json:"profiles"`
}

type GetChapterVideoProfile struct {
	Resolution   string `json:"resolution"`
	PlaylistPath string `json:"playlist_path"`
	FileSize     int64  `json:"file_size"`
}
