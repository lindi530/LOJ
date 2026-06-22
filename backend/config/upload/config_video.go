package upload

// Video holds video upload and transcode settings.
type Video struct {
	Path                         string `mapstructure:"path"`                             // Video origin path saved for business records.
	ChunkSize                    int64  `mapstructure:"chunk_size"`                       // Recommended upload chunk size returned to clients.
	MaxPartCount                 int    `mapstructure:"max_part_count"`                   // Maximum multipart upload part count.
	ChunkUploadURLExpiresMinutes int    `mapstructure:"chunk_upload_url_expires_minutes"` // Presigned chunk upload URL expiration in minutes.
	FFmpegPath                   string `mapstructure:"ffmpeg_path"`                      // ffmpeg executable path.
	FFprobePath                  string `mapstructure:"ffprobe_path"`                     // ffprobe executable path.
	Cover                        Cover  `mapstructure:"cover"`                            // Video cover extraction settings.
	HLS                          HLS    `mapstructure:"hls"`                              // HLS transcode settings.
}

// Cover holds video cover extraction settings.
type Cover struct {
	FrameCount int `mapstructure:"frame_count"` // Number of frames to capture for the cover.
	Quality    int `mapstructure:"quality"`     // ffmpeg image quality value for the cover.
}

// HLS holds HLS transcode settings.
type HLS struct {
	PlaylistName     string       `mapstructure:"playlist_name"`      // Master playlist file name.
	IndexName        string       `mapstructure:"index_name"`         // Variant playlist file name.
	SegmentSeconds   int          `mapstructure:"segment_seconds"`    // HLS segment duration in seconds.
	AudioBitrateKbps int          `mapstructure:"audio_bitrate_kbps"` // Audio bitrate in kbps.
	Preset           string       `mapstructure:"preset"`             // x264 encoding preset.
	Profile          string       `mapstructure:"profile"`            // H.264 profile.
	Variants         []HLSVariant `mapstructure:"variants"`           // Output resolution variants.
}

// HLSVariant describes one HLS output profile.
type HLSVariant struct {
	Dir         string `mapstructure:"dir"`          // Output directory name for this variant.
	Resolution  string `mapstructure:"resolution"`   // Display resolution label.
	Height      int    `mapstructure:"height"`       // Target video height.
	BitrateKbps int    `mapstructure:"bitrate_kbps"` // Target video bitrate in kbps.
}
