package storage

import (
	"bytes"
	"context"
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/enums"
	"mime/multipart"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client  *s3.Client
	buckets map[enums.BucketType]string
}

func NewS3Storage() *S3Storage {
	buckets := make(map[enums.BucketType]string)
	buckets[enums.PetProfilePic] = bootstrap.Run().Env.Storage.Buckets.PetProfilePic
	buckets[enums.UserProfilePic] = bootstrap.Run().Env.Storage.Buckets.UserProfilePic
	buckets[enums.PetSitterFile] = bootstrap.Run().Env.Storage.Buckets.PetSitterFile
	buckets[enums.ChatMedia] = bootstrap.Run().Env.Storage.Buckets.ChatMedia
	return &S3Storage{
		buckets: buckets,
	}
}

func (s3Storage *S3Storage) setS3Client(bucketType enums.BucketType) error {
	bucketTypes := enums.GetAllBucketTypes()
	if !slices.Contains(bucketTypes, bucketType) {
		return fmt.Errorf("bucket not exist")
	}

	if s3Storage.client != nil {
		return nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("default"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			bootstrap.Run().Env.Storage.AccessKey,
			bootstrap.Run().Env.Storage.SecretKey,
			"",
		)),
	)
	if err != nil {
		return err
	}

	s3Storage.client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(bootstrap.Run().Env.Storage.Endpoint)
	})
	return nil
}

func (s3Storage *S3Storage) UploadFile(bucketType enums.BucketType, key string, file *multipart.FileHeader) error {
	err := s3Storage.setS3Client(bucketType)
	if err != nil {
		return err
	}
	bucket := s3Storage.buckets[bucketType]

	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	_, err = s3Storage.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fileReader,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s3Storage *S3Storage) UploadBytes(bucketType enums.BucketType, key string, contentType string, data []byte) error {
	err := s3Storage.setS3Client(bucketType)
	if err != nil {
		return err
	}
	bucket := s3Storage.buckets[bucketType]

	_, err = s3Storage.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s3Storage *S3Storage) DeleteObject(bucketType enums.BucketType, key string) error {
	err := s3Storage.setS3Client(bucketType)
	if err != nil {
		return err
	}
	bucket := s3Storage.buckets[bucketType]

	_, err = s3Storage.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s3Storage *S3Storage) GetPresignedURL(bucketType enums.BucketType, key string, expiration time.Duration) (string, error) {
	err := s3Storage.setS3Client(bucketType)
	if err != nil {
		return "", err
	}
	bucket := s3Storage.buckets[bucketType]

	req := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	presigner := s3.NewPresignClient(s3Storage.client)

	presignedURL, err := presigner.PresignGetObject(context.TODO(), req, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", err
	}

	return presignedURL.URL, nil
}
