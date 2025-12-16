package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// TODO: find secure but efficient timeframe for duration
const (
	UPLOAD_TIME   time.Duration = 24 * time.Hour
	DOWNLOAD_TIME time.Duration = 24 * time.Hour
)

type Storage struct {
	PresignClient *s3.PresignClient
}

func InitializeStorage() *Storage {
	// load the config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading config file")
	}

	// create s3 clients
	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	// return wrapped s3 object
	return &Storage{
		PresignClient: presignClient,
	}
}

func CreatePictureId() string {
	return uuid.New().String()
}

func (s *Storage) GenerateUploadURL(bucket string, key string) (string, error) {
	req := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	url, err := s.PresignClient.PresignPutObject(
		context.TODO(),
		req,
		s3.WithPresignExpires(UPLOAD_TIME),
	)

	if err != nil {
		return "", err
	}
	return url.URL, nil
}

func (s *Storage) GenerateDownloadURL(bucket string, key string) (string, error) {
	// TODO: right now assuming image/jpeg but analyze if this is best decision
	req := &s3.GetObjectInput{
		Bucket:                     &bucket,
		Key:                        &key,
		ResponseContentType:        aws.String("image/jpeg"),
		ResponseContentDisposition: aws.String("inline"),
	}

	url, err := s.PresignClient.PresignGetObject(
		context.TODO(),
		req,
		s3.WithPresignExpires(DOWNLOAD_TIME),
	)

	if err != nil {
		return "", err
	}

	return url.URL, nil
}
