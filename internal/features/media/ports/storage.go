package ports

import "time"

type ObjectStorage interface {
	GeneratePresignedPutURL(key string, contentType string, ttl time.Duration) (string, error)
}
