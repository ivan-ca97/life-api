package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/ports"
)

type r2Lister struct {
	client *s3.Client
	bucket string
}

var _ ports.ObjectLister = (*r2Lister)(nil)
var _ ports.ObjectDeleter = (*r2Lister)(nil)

func NewR2Lister(accountId, accessKeyId, secretAccessKey, bucket string) *r2Lister {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)
	options := s3.Options{
		Region:       "auto",
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, "")),
		BaseEndpoint: aws.String(endpoint),
	}
	client := s3.New(options)
	return &r2Lister{
		client: client,
		bucket: bucket,
	}
}

// DeleteKeys removes up to 1000 keys per S3 batch-delete request.
func (l *r2Lister) DeleteKeys(keys []string) error {
	const batchSize = 1000
	for i := 0; i < len(keys); i += batchSize {
		end := i + batchSize
		if end > len(keys) {
			end = len(keys)
		}
		batch := keys[i:end]
		objects := make([]types.ObjectIdentifier, len(batch))
		for j, k := range batch {
			key := k
			objects[j] = types.ObjectIdentifier{Key: aws.String(key)}
		}
		out, err := l.client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
			Bucket: aws.String(l.bucket),
			Delete: &types.Delete{Objects: objects, Quiet: aws.Bool(true)},
		})
		if err != nil {
			return fmt.Errorf("batch delete R2 objects: %w", err)
		}
		if len(out.Errors) > 0 {
			return fmt.Errorf("batch delete R2: %d errors, first: %s %s", len(out.Errors), aws.ToString(out.Errors[0].Key), aws.ToString(out.Errors[0].Message))
		}
	}
	return nil
}

func (l *r2Lister) ListAllKeys(prefix string) ([]string, error) {
	var keys []string
	var continuationToken *string
	for {
		out, err := l.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(l.bucket),
			Prefix:            aws.String(prefix),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, fmt.Errorf("listing R2 objects: %w", err)
		}
		for _, obj := range out.Contents {
			if obj.Key != nil {
				keys = append(keys, *obj.Key)
			}
		}
		if !*out.IsTruncated || out.NextContinuationToken == nil {
			break
		}
		continuationToken = out.NextContinuationToken
	}
	return keys, nil
}
