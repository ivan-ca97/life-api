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
	UploadUrl string
	PublicUrl string
}

type MediaService interface {
	GenerateUploadUrl(request UploadRequest) (*UploadResult, error)
}

type AuthorizedMediaService interface {
	GenerateUploadUrl(ctx context.Context, request UploadRequest) (*UploadResult, error)
}
