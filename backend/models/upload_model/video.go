package upload_model

type CreateVideoUploadTaskReq struct{}

type CreateVideoUploadTaskResp struct {
	UploadID      int64  `json:"upload_id"`
	MinioUploadID string `json:"minio_upload_id"`
	ChunkSize     int    `json:"chunk_size"`
}

type CreateChunkUploadURLReq struct {
	UploadID int64  `json:"upload_id"`
	ChunkID  int    `json:"chunk_id"`
	MD5      string `json:"md5"`
}

type CreateChunkUploadURLResp struct {
	Exist      bool              `json:"exist"`
	Method     string            `json:"method,omitempty"`
	UploadURL  string            `json:"upload_url,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	PartNumber int               `json:"part_number,omitempty"`
}

type VideoFinishReq struct {
	UploadID   int64  `json:"upload_id"`
	MD5        string `json:"md5"`
	FileName   string `json:"file_name"`
	ChunkCount int    `json:"chunk_count"`
}

type VideoFinishResp struct {
	OriginPath string `json:"origin_path"`
	ID         int64  `json:"id"`
}
