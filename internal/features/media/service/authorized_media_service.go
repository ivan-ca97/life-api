package service

import (
	"context"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/media/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedMediaService struct {
	base       ports.MediaService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedMediaService = (*authorizedMediaService)(nil)

func NewAuthorizedMediaService(base ports.MediaService, authorizer auth.AuthorizationService) *authorizedMediaService {
	return &authorizedMediaService{
		base:       base,
		authorizer: authorizer,
	}
}

// GenerateUploadUrl issues an upload URL inside the owner's storage space, so the
// actor must be allowed to create media for that user.
func (s *authorizedMediaService) GenerateUploadUrl(ctx context.Context, request ports.UploadRequest) (*ports.UploadResult, error) {
	err := s.authorizer.Authorize(ctx, request.UserId, permissions.MediaCreate)
	if err != nil {
		return nil, err
	}
	result, err := s.base.GenerateUploadUrl(request)
	if err != nil {
		return nil, err
	}
	return result, nil
}
