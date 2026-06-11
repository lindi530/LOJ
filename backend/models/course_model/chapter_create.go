package course_model

type ChapterCreateReq struct {
	CourseID    int64     `json:"course_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Sort        int       `json:"sort"`
	Problems    []problem `json:"problems"`
	VideoURL    string    `json:"video_url"`
}

type problem struct {
	ID   int64 `json:"problem_id"`
	Sort int   `json:"sort"`
}

type ChapterCreateResp struct {
}

/*
{
  "course_id": 1,
  "title": "章节标题",
  "description": null,
  "sort": 1, // 前端上传题目的顺序
  "problems": [
    {
      "problem_id": 101,
      "sort": 1
    }
  ],
  "video_url": "https://example.com/video.mp4"
}
*/
