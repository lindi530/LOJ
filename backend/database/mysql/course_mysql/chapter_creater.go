package course_mysql

import (
	"GO1/global"
	"GO1/models/course_model"

	"gorm.io/gorm"
)

func CreateChapter(
	chapter *course_model.CourseChapter,
	problems []course_model.ChapterProblem,
	video *course_model.ChapterVideo,
) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("course_chapters").Create(chapter).Error; err != nil {
			return err
		}

		if len(problems) > 0 {
			for i := range problems {
				problems[i].ChapterID = chapter.ID
			}

			if err := tx.Table("chapter_problems").Create(&problems).Error; err != nil {
				return err
			}
		}

		if video != nil {
			video.ChapterID = chapter.ID
			if err := tx.Table("chapter_videos").Create(video).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
