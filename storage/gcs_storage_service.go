package storage

import (
	"io"
	"os"
	"errors"
	"bufio"
	"context"
	"encoding/base64"
	"database-backup-manager/constants"
	"cloud.google.com/go/storage"
    "google.golang.org/api/option"
)

type gcsStorageService struct {
	client *storage.Client
	bucket string
}

func NewGcsStorageService() (*gcsStorageService,error) {
	secretKey := os.Getenv(constants.STORAGE_SECRETKEY)
	if (len(secretKey) <= 0) {
		return nil, errors.New("No credentials in secretKey specified!")
	}

	bucket := os.Getenv(constants.STORAGE_BUCKET)
	if (len(bucket) <= 0) {
		return nil, errors.New("No bucket specified!")
	}

	jsonCred, err := base64.StdEncoding.DecodeString(secretKey)
    if err != nil {
		panic(err.Error())
    }
	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON(jsonCred))
	if (err != nil) {
		return nil, err
	}

	ss := new(gcsStorageService)
	ss.client = client
	ss.bucket = bucket
	return ss, nil
}

func (ss *gcsStorageService) UploadFile(localSourceFile string, remoteTargetFile string, mimeType string) error {
	wc := ss.client.Bucket(ss.bucket).Object(remoteTargetFile).NewWriter(context.Background())
    wc.ContentType = mimeType
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	defer wc.Close()
	
	file, err := os.Open(localSourceFile)
	if (err != nil) {return err}
	defer file.Close()
	reader := bufio.NewReader(file)
	bufferSize := 4 * 1024;
    buffer := make([]byte, bufferSize)

	for {
        numberOfBytes, err := reader.Read(buffer)
        if err != nil {
            if err != io.EOF {
                return err
            }
            break
		}

		if (numberOfBytes < bufferSize) {
			buffer = buffer[0:numberOfBytes]
		}

		if _, err := wc.Write(buffer); err != nil {
			return err
		}
    }

	return nil
}