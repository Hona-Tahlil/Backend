package domainstorage

import (
	"hona/backend/internal/domain/enums"
	"mime/multipart"
	"time"
)

type Storage interface {
	UploadFile(bucketType enums.BucketType, key string, file *multipart.FileHeader) error
	DeleteObject(bucketType enums.BucketType, key string) error
	GetPresignedURL(bucketType enums.BucketType, key string, expiration time.Duration) (string, error)
}

