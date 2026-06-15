package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/ports"
)

type r2Lister struct {
	client *s3.Client
	bucket string
}

var _ ports.ObjectLister = (*r2Lister)(nil)

func NewR2Lister(accountId, accessKeyId, secretAccessKey, bucket string) *r2Lister {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)
	client := s3.New(s3.Options{
		Region:       "auto",
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, "")),
		BaseEndpoint: aws.String(endpoint),
	})
	return &r2Lister{client: client, bucket: bucket}
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
