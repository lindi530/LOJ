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
	ID        int64     `json:"id"`
	Title     *string   `json:"title"`
	URL       string    `json:"url"`
	Duration  int       `json:"duration"`
	SizeBytes *int64    `json:"sizeBytes"`
	CreatedBy *int64    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type CourseLesson struct {
	ID            int64     `json:"id"`
	CourseID      int64     `json:"courseId"`
	ChapterID     int64     `json:"chapterId"`
	VideoAssetID  *int64    `json:"videoAssetId"`
	Title         string    `json:"title"`
	Duration      int       `json:"duration"`
	IsFreePreview int8      `json:"isFreePreview"`
	Sort          int       `json:"sort"`
	Status        int8      `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type LessonProblem struct {
	ID          int64     `json:"id"`
	CourseID    int64     `json:"courseId"`
	LessonID    int64     `json:"lessonId"`
	ProblemID   int64     `json:"problemId"`
	ProblemType int8      `json:"problemType"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LessonProgress struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"userId"`
	CourseID        int64      `json:"courseId"`
	LessonID        int64      `json:"lessonId"`
	LastPosition    int        `json:"lastPosition"`
	WatchedSeconds  int        `json:"watchedSeconds"`
	ProgressPercent float64    `json:"progressPercent"`
	IsFinished      int8       `json:"isFinished"`
	FinishedAt      *time.Time `json:"finishedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
