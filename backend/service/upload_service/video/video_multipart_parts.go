package video

import (
	"context"
	"strings"

	"GO1/database/redis/course_redis"
	"GO1/global"

	miniogo "github.com/minio/minio-go/v7"
)

func getUploadedVideoPart(ctx context.Context, meta course_redis.VideoMultipartUploadMeta, partNumber int) (miniogo.ObjectPart, bool, error) {
	result, err := global.MinioCore.ListObjectParts(ctx, meta.Bucket, meta.TempObject, meta.MinioUploadID, partNumber-1, 1)
	if err != nil {
		return miniogo.ObjectPart{}, false, err
	}
	for _, part := range result.ObjectParts {
		if part.PartNumber == partNumber {
			return part, true, nil
		}
	}
	return miniogo.ObjectPart{}, false, nil
}

func listUploadedVideoParts(ctx context.Context, meta course_redis.VideoMultipartUploadMeta) (map[int]miniogo.ObjectPart, error) {
	parts := make(map[int]miniogo.ObjectPart)
	partNumberMarker := 0

	for {
		result, err := global.MinioCore.ListObjectParts(ctx, meta.Bucket, meta.TempObject, meta.MinioUploadID, partNumberMarker, 1000)
		if err != nil {
			return nil, err
		}
		for _, part := range result.ObjectParts {
			parts[part.PartNumber] = part
		}
		if !result.IsTruncated {
			break
		}
		partNumberMarker = result.NextPartNumberMarker
	}

	return parts, nil
}

func objectPartToCompletePart(part miniogo.ObjectPart) miniogo.CompletePart {
	return miniogo.CompletePart{
		PartNumber:        part.PartNumber,
		ETag:              part.ETag,
		ChecksumCRC32:     part.ChecksumCRC32,
		ChecksumCRC32C:    part.ChecksumCRC32C,
		ChecksumSHA1:      part.ChecksumSHA1,
		ChecksumSHA256:    part.ChecksumSHA256,
		ChecksumCRC64NVME: part.ChecksumCRC64NVME,
		ChecksumMD5:       part.ChecksumMD5,
		ChecksumSHA512:    part.ChecksumSHA512,
		ChecksumXXHash64:  part.ChecksumXXHash64,
		ChecksumXXHash3:   part.ChecksumXXHash3,
		ChecksumXXHash128: part.ChecksumXXHash128,
	}
}

func objectPartMatchesMD5(part miniogo.ObjectPart, md5 string) bool {
	return md5 == "" || strings.EqualFold(strings.Trim(part.ETag, "\""), md5)
}
