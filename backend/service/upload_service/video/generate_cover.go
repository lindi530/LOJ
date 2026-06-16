package video

import "context"

func generateCover(ctx context.Context, inputPath string, coverPath string) error {
	return runFFmpeg(
		ctx,
		"-y",
		"-i", inputPath,
		"-frames:v", "1",
		"-q:v", "2",
		coverPath,
	)
}
