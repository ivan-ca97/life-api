package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedDailyPhotoService struct {
	base       ports.PhotoService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedPhotoService = (*authorizedDailyPhotoService)(nil)

func NewAuthorizedDailyPhotoService(base ports.PhotoService, authorizer auth.AuthorizationService) *authorizedDailyPhotoService {
	return &authorizedDailyPhotoService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedDailyPhotoService) Create(ctx context.Context, ownerId uuid.UUID, params ports.CreatePhotoParams) (*domain.DailyPhoto, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate)
	if err != nil {
		return nil, err
	}
	photo, err := s.base.Create(ownerId, params)
	if err != nil {
		return nil, err
	}
	return photo, nil
}

func (s *authorizedDailyPhotoService) List(ctx context.Context, ownerId uuid.UUID, date time.Time) ([]domain.DailyPhoto, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyRead)
	if err != nil {
		return nil, err
	}
	photos, err := s.base.List(ownerId, date)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func (s *authorizedDailyPhotoService) Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params ports.UpdatePhotoParams) (*domain.DailyPhoto, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate)
	if err != nil {
		return nil, err
	}
	photo, err := s.base.Update(id, ownerId, params)
	if err != nil {
		return nil, err
	}
	return photo, nil
}

func (s *authorizedDailyPhotoService) Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate)
	if err != nil {
		return err
	}
	err = s.base.Delete(id, ownerId)
	if err != nil {
		return err
	}
	return nil
}
