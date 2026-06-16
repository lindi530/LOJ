package minio

import (
	"os"

	"GO1/global"

	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() {
	address := os.Getenv("MINIO_ADDRESS")
	if address == "" {
		address = global.Config.Minio.Address
	}
	region := os.Getenv("MINIO_REGION")
	if region == "" {
		region = global.Config.Minio.Region
	}
	if region == "" {
		region = "us-east-1"
	}

	client, err := newMinioCore(address, region)
	if err != nil {
		global.Logger.Fatalf("minio init failed: %v", err)
	}

	publicAddress := os.Getenv("MINIO_PUBLIC_ADDRESS")
	if publicAddress == "" {
		publicAddress = global.Config.Minio.PublicAddress
	}
	if publicAddress == "" {
		publicAddress = address
	}

	presignClient, err := newMinioCore(publicAddress, region)
	if err != nil {
		global.Logger.Fatalf("minio presign init failed: %v", err)
	}

	global.MinioCore = client
	global.MinioPresignCore = presignClient
}

func newMinioCore(address string, region string) (*miniogo.Core, error) {
	return miniogo.NewCore(address, &miniogo.Options{
		Creds:        credentials.NewStaticV4(global.Config.Minio.AccessKeyID, global.Config.Minio.SecretAccessKey, ""),
		Secure:       false, // https
		Region:       region,
		BucketLookup: miniogo.BucketLookupPath,
	})
}
