package video

import (
	"context"
	"strconv"

	"GO1/global"
)

func generateCover(ctx context.Context, inputPath string, coverPath string) error {
	return runFFmpeg(
		ctx,
		"-y",
		"-i", inputPath,
		"-frames:v", strconv.Itoa(global.Config.Upload.Video.Cover.FrameCount),
		"-q:v", strconv.Itoa(global.Config.Upload.Video.Cover.Quality),
		coverPath,
	)
}
