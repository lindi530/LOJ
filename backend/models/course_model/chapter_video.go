package course_model

type GetChapterVideoReq struct {
	CourseID  int64 `uri:"course_id" binding:"required"`
	ChapterID int64 `uri:"chapter_id" binding:"required"`
	VideoID   int64 `uri:"video_id" binding:"required"`
}

type GetVideoHLSKeyReq struct {
	VideoID int64  `uri:"video_id"`
	Variant string `form:"variant"`
}

type GetChapterVideoResp struct {
	ID          int64                    `json:"id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	CoverURL    string                   `json:"cover_url"`
	Duration    int                      `json:"duration"`
	SizeBytes   int64                    `json:"size_bytes"`
	Profiles    []GetChapterVideoProfile `json:"profiles"`
}

type GetChapterVideoProfile struct {
	Resolution string `json:"resolution"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Bitrate    int    `json:"bitrate"`
	PlayURL    string `json:"play_url"`
	FileSize   int64  `json:"file_size"`
}
