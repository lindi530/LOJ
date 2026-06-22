package video

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"GO1/global"
)

func runFFmpeg(ctx context.Context, args ...string) error {
	return runFFmpegInDir(ctx, "", args...)
}

func runFFmpegInDir(ctx context.Context, workDir string, args ...string) error {
	cmd := exec.CommandContext(ctx, global.Config.Upload.Video.FFmpegPath, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
