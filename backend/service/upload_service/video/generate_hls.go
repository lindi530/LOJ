package video

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"GO1/config/upload"
	"GO1/global"
	"GO1/models/course_model"
)

func generateHLS(ctx context.Context, inputPath string, outputDir string) (string, []course_model.VideoProfile, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", nil, err
	}

	hlsConfig := global.Config.Upload.Video.HLS
	profiles := make([]course_model.VideoProfile, 0, len(hlsConfig.Variants))
	for _, variant := range hlsConfig.Variants {
		variantDir := filepath.Join(outputDir, variant.Dir)
		if err := generateHLSVariant(ctx, inputPath, variantDir, variant, hlsConfig); err != nil {
			return "", nil, err
		}

		segmentPath, err := firstHLSSegmentPath(variantDir)
		if err != nil {
			return "", nil, err
		}
		width, height, err := probeVideoDimensions(ctx, segmentPath)
		if err != nil {
			return "", nil, err
		}

		fileSize, err := directoryFileSize(variantDir)
		if err != nil {
			return "", nil, err
		}

		profiles = append(profiles, course_model.VideoProfile{
			Resolution:   variant.Resolution,
			Width:        width,
			Height:       height,
			Bitrate:      variant.BitrateKbps,
			PlaylistPath: strings.ReplaceAll(filepath.ToSlash(filepath.Join(variant.Dir, hlsConfig.IndexName)), "\\", "/"),
			FileSize:     fileSize,
		})
	}

	if err := writeMasterPlaylist(outputDir, profiles, hlsConfig); err != nil {
		return "", nil, err
	}

	return filepath.Join(outputDir, hlsConfig.PlaylistName), profiles, nil
}

func generateHLSVariant(ctx context.Context, inputPath string, outputDir string, variant upload.HLSVariant, hlsConfig upload.HLS) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	maxrateKbps := variant.BitrateKbps * 12 / 10
	bufsizeKbps := variant.BitrateKbps * 2
	segmentSeconds := strconv.Itoa(hlsConfig.SegmentSeconds)

	return runFFmpegInDir(
		ctx,
		outputDir,
		"-y",
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=-2:%d", variant.Height),
		"-c:v", "libx264",
		"-preset", hlsConfig.Preset,
		"-profile:v", hlsConfig.Profile,
		"-b:v", fmt.Sprintf("%dk", variant.BitrateKbps),
		"-maxrate", fmt.Sprintf("%dk", maxrateKbps),
		"-bufsize", fmt.Sprintf("%dk", bufsizeKbps),
		"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%s)", segmentSeconds),
		"-c:a", "aac",
		"-b:a", fmt.Sprintf("%dk", hlsConfig.AudioBitrateKbps),
		"-ac", "2",
		"-hls_time", segmentSeconds,
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-hls_segment_filename", "%03d.ts",
		hlsConfig.IndexName,
	)
}

func writeMasterPlaylist(outputDir string, profiles []course_model.VideoProfile, hlsConfig upload.HLS) error {
	var builder strings.Builder
	builder.WriteString("#EXTM3U\n")
	builder.WriteString("#EXT-X-VERSION:6\n")
	for _, profile := range profiles {
		bandwidth := (profile.Bitrate + hlsConfig.AudioBitrateKbps) * 1000
		fmt.Fprintf(
			&builder,
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n%s\n",
			bandwidth,
			profile.Width,
			profile.Height,
			profile.PlaylistPath,
		)
	}

	return os.WriteFile(filepath.Join(outputDir, hlsConfig.PlaylistName), []byte(builder.String()), 0644)
}

func firstHLSSegmentPath(outputDir string) (string, error) {
	segments, err := filepath.Glob(filepath.Join(outputDir, "*.ts"))
	if err != nil {
		return "", err
	}
	if len(segments) == 0 {
		return "", fmt.Errorf("hls segment not found in %s", outputDir)
	}
	sort.Strings(segments)
	return segments[0], nil
}

func directoryFileSize(dir string) (int64, error) {
	var size int64
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})
	return size, err
}

func probeVideoDimensions(ctx context.Context, filePath string) (int, int, error) {
	cmd := exec.CommandContext(
		ctx,
		global.Config.Upload.Video.FFprobePath,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "json",
		filePath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	var result struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, 0, err
	}
	if len(result.Streams) == 0 || result.Streams[0].Width <= 0 || result.Streams[0].Height <= 0 {
		return 0, 0, fmt.Errorf("video stream dimensions not found: %s", filePath)
	}

	return result.Streams[0].Width, result.Streams[0].Height, nil
}
