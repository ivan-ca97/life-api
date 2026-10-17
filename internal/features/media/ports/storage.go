package ports

import "time"

type ObjectStorage interface {
	GeneratePresignedPutUrl(key string, contentType string, ttl time.Duration) (string, error)
}
