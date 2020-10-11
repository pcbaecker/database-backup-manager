package storage

import (
	"os"
	"context"
	"database-backup-manager/constants"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type s3StorageService struct {
	client *minio.Client
	bucket string
}

func NewS3StorageService() (*s3StorageService,error) {
	endpoint := os.Getenv(constants.STORAGE_ENDPOINT)
	accessKey := os.Getenv(constants.STORAGE_ACCESSKEY)
	secretKey := os.Getenv(constants.STORAGE_SECRETKEY)
	bucket := os.Getenv(constants.STORAGE_BUCKET)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	
	ss := new(s3StorageService)
	ss.client = minioClient
	ss.bucket = bucket
	return ss, nil
}

func (ss *s3StorageService) UploadFile(localSourceFile string, remoteTargetFile string, mimeType string) error {
	_, err := ss.client.FPutObject(
		context.Background(), 
		ss.bucket, 
		remoteTargetFile, 
		localSourceFile, 
		minio.PutObjectOptions{ContentType: mimeType})
	return err
}