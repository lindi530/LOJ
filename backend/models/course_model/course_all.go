package course_model

import "time"

type Course struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	CoverURL     *string   `json:"coverUrl"`
	Price        float64   `json:"price"`
	IsFree       int8      `json:"isFree"`
	LessonCount  int       `json:"lessonCount"`
	StudentCount int       `json:"studentCount"`
	Status       int8      `json:"status"`
	Sort         int       `json:"sort"`
	CreatedBy    *int64    `json:"createdBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CourseChapter struct {
	ID          int64     `json:"id"`
	CourseID    int64     `json:"courseId"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type VideoAsset struct {
	ID         int64  `json:"id"`
	OriginPath string `json:"origin_path"`
	PlayPath   string `json:"play_path"`
	CoverPath  string `json:"cover_path"`
	Status     int8   `json:"status"`
	Duration   int    `json:"duration"`
	SizeBytes  *int64 `json:"sizeBytes"`
}

type VideoProfile struct {
	ID           int64     `json:"id"`
	VideoID      int64     `json:"video_id"`   // 对应 原视频VideoAsset 中的ID
	Resolution   string    `json:"resolution"` // 480P, 720P, 1080P
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	Bitrate      int       `json:"bitrate"`       // 码率(kbps)
	PlaylistPath string    `json:"playlist_path"` // index.m3u8
	FileSize     int64     `json:"file_size"`     // 字节
	CreatedAt    time.Time `json:"created_at"`
}

const (
	VideoAssetStatusUploading int8 = iota
	VideoAssetStatusTranscoding
	VideoAssetStatusPlayable
	VideoAssetStatusTranscodeFailed
)

//type CourseLesson struct {
//	ID            int64     `json:"id"`
//	CourseID      int64     `json:"courseId"`
//	ChapterID     int64     `json:"chapterId"`
//	VideoAssetID  *int64    `json:"videoAssetId"`
//	Title         string    `json:"title"`
//	Duration      int       `json:"duration"`
//	IsFreePreview int8      `json:"isFreePreview"`
//	Sort          int       `json:"sort"`
//	Status        int8      `json:"status"`
//	CreatedAt     time.Time `json:"createdAt"`
//	UpdatedAt     time.Time `json:"updatedAt"`
//}

type ChapterProblem struct {
	ID        int64 `json:"id"`
	CourseID  int64 `json:"courseId"`
	ChapterID int64 `json:"chapterId"`
	ProblemID int64 `json:"problemId"`
	Sort      int   `json:"sort"`
	//CreatedAt   time.Time `json:"createdAt"`
	//ProblemType int8      `json:"problemType"`
	//LessonID    int64     `json:"lessonId"`
}

type ChapterVideo struct {
	ID        int64 `json:"id"`
	CourseID  int64 `json:"courseId"`
	ChapterID int64 `json:"chapterId"`
	VideoID   int64 `json:"videoId"`
	//ProblemType int8      `json:"problemType"`
	//CreatedAt   time.Time `json:"createdAt"`
	//LessonID    int64     `json:"lessonId"`
}

//type LessonProgress struct {
//	ID              int64      `json:"id"`
//	UserID          int64      `json:"userId"`
//	CourseID        int64      `json:"courseId"`
//	LessonID        int64      `json:"lessonId"`
//	LastPosition    int        `json:"lastPosition"`
//	WatchedSeconds  int        `json:"watchedSeconds"`
//	ProgressPercent float64    `json:"progressPercent"`
//	IsFinished      int8       `json:"isFinished"`
//	FinishedAt      *time.Time `json:"finishedAt"`
//	CreatedAt       time.Time  `json:"createdAt"`
//	UpdatedAt       time.Time  `json:"updatedAt"`
//}
