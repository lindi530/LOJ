package upload

type Video struct {
	Path      string `mapstructure:"path"`
	ChunkSize int64  `mapstructure:"chunk_size"`
}
