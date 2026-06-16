package course_model

type ChapterCreateReq struct {
	CourseID    int64                  `json:"course_id" binding:"required"`
	Title       string                 `json:"title" binding:"required"`
	Description *string                `json:"description"`
	Sort        int                    `json:"sort" binding:"required"`
	Problems    []ChapterCreateProblem `json:"problems"`
	Video       ChapterCreateVideo     `json:"video" binding:"required"`
}

type ChapterCreateProblem struct {
	ID   int64 `json:"id"`
	Sort int   `json:"sort"`
}

type ChapterCreateVideo struct {
	OriginPath string `json:"origin_path"`
	ID         int64  `json:"id"`
}

type ChapterCreateResp struct {
}

/*
{
  "course_id": 1,
  "title": "chapter title",
  "description": null,
  "sort": 1,
  "problems": [
    {
      "id": 101,
      "sort": 1
    }
  ],
  "video": {
    "origin_path": "video/10/example.mp4",
    "id": 10
  }
}
*/
