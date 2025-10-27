package storage

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/enums"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/credentials"
	"github.com/aws/aws-sdk-go-v2/aws/session"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client  *s3.Client
	buckets map[enums.BucketType]string
}

func NewS3Storage() *S3Storage {
	buckets := make(map[enums.BucketType]string)
	buckets[enums.PetProfilePic] = bootstrap.Run().Env.Storage.Buckets.PetProfilePic
	return &S3Storage{
		buckets: buckets,
	}
}

func (s3StorageS3Storage *S3Storage) setS3Client(bucketType enum.BucketType) error {
	bucketTypes := enums.GetAllBucketTypes()
	if !slices.Contains(bucketTypes, bucketType) {
		return fmt.Errorf("bucket not exist")
	}
	if s3StorageS3Storage.client != nil {
		return nil
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s3StorageS3Storage.storage.AccessKey, s3StorageS3Storage.storage.SecretKey, ""),
		Region:      aws.String(s3StorageS3Storage.storage.Region),
		Endpoint:    aws.String(s3StorageS3Storage.storage.Endpoint),
	})

	if err != nil {
		return fmt.Errorf("unable to create AWS session, %w", err)
	}

	s3StorageS3Storage.clients = s3.New(sess)
	return nil
}
