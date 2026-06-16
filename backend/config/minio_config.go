package config

type Minio struct {
	Address         string `mapstructure:"address"`
	PublicAddress   string `mapstructure:"public_address"`
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          Bucket `mapstructure:"bucket"`
}

type Bucket struct {
	Video          string `mapstructure:"video"`
	VideoPath      string `mapstructure:"video_path"`
	VideoPlayPath  string `mapstructure:"video_play_path"`
	VideoCoverPath string `mapstructure:"video_cover_path"`
}

/*
minio:
    address: "minio:9000"
    public_address: "localhost:9000"
    region: "us-east-1"
    access_key_id: "minioadmin"
    secret_access_key: "minioadmin"
    bucket:
        video: "loj"
        video_path: "raw"
        video_play_path: "hls"
        video_cover_path: "cover"
*/
