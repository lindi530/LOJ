package upload_model

import "mime/multipart"

type CreateVideoUploadTaskReq struct{}

type CreateVideoUploadTaskResp struct {
	UploadID  int64 `json:"upload_id"`
	ChunkSize int   `json:"chunk_size"`
}
type CheckChunkMD5Req struct {
	UploadID int64  `json:"upload_id"`
	ChunkID  int    `json:"chunk_id"`
	MD5      string `json:"md5"`
}

type CheckChunkMD5Resp struct {
	Exist bool `json:"exist"`
}

type ReceiveChunkReq struct {
	UploadID int64                 `json:"upload_id" form:"upload_id"`
	ChunkID  int                   `json:"chunk_id" form:"chunk_id"`
	Filename string                `json:"filename" form:"filename"`
	Chunk    *multipart.FileHeader `json:"-" form:"chunk"`
}

type ReceiveChunkResp struct {
}

type VideoFinishReq struct {
	UploadID   int64  `json:"upload_id"`
	MD5        string `json:"md5"`
	FileName   string `json:"file_name"`
	ChunkCount int    `json:"chunk_count"`
}

type VideoFinishResp struct {
	URL string `json:"url"`
	ID  int64  `json:"id"`
}
