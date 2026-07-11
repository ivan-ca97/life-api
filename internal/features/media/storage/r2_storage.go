package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ivan-ca97/life/internal/features/media/ports"
)

type r2Storage struct {
	presignClient *s3.PresignClient
	bucket        string
}

var _ ports.ObjectStorage = (*r2Storage)(nil)

func NewR2Storage(accountId, accessKeyId, secretAccessKey, bucket string) *r2Storage {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)

	options := s3.Options{
		Region:       "auto",
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, "")),
		BaseEndpoint: aws.String(endpoint),
	}

	client := s3.New(options)

	return &r2Storage{
		presignClient: s3.NewPresignClient(client),
		bucket:        bucket,
	}
}

func (s *r2Storage) GeneratePresignedPutURL(key string, contentType string, ttl time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	presigned, err := s.presignClient.PresignPutObject(context.Background(), input,
		s3.WithPresignExpires(ttl),
	)
	if err != nil {
		return "", fmt.Errorf("generating presigned URL: %w", err)
	}

	return presigned.URL, nil
}
