package ports

import (
	"context"

	"github.com/google/uuid"
)

type UploadRequest struct {
	UserId      uuid.UUID
	Filename    string
	ContentType string
}

type UploadResult struct {
	UploadURL string
	PublicURL string
}

type MediaService interface {
	GenerateUploadURL(ctx context.Context, request UploadRequest) (*UploadResult, error)
}
