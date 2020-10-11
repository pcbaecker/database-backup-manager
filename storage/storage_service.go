package storage

import (
	"os"
	"errors"
	"database-backup-manager/constants"
)

type StorageService interface {
	UploadFile(string, string , string) error
}

func NewStorageService() (StorageService, error) {
	storage_type := os.Getenv(constants.STORAGE_TYPE)
	if (storage_type == "gcs") {
		return NewGcsStorageService()
	}
	if (storage_type == "s3") {
		return NewS3StorageService()
	}
	return nil, errors.New("Unspecified storage type = " + storage_type)
}