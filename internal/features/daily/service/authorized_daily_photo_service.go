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
	return &authorizedDailyPhotoService{base: base, authorizer: authorizer}
}

func (s *authorizedDailyPhotoService) Create(ctx context.Context, ownerId uuid.UUID, params ports.CreatePhotoParams) (*domain.DailyPhoto, error) {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate); err != nil {
		return nil, err
	}
	return s.base.Create(ownerId, params)
}

func (s *authorizedDailyPhotoService) List(ctx context.Context, ownerId uuid.UUID, date time.Time) ([]domain.DailyPhoto, error) {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyRead); err != nil {
		return nil, err
	}
	return s.base.List(ownerId, date)
}

func (s *authorizedDailyPhotoService) Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params ports.UpdatePhotoParams) (*domain.DailyPhoto, error) {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate); err != nil {
		return nil, err
	}
	return s.base.Update(id, ownerId, params)
}

func (s *authorizedDailyPhotoService) Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate); err != nil {
		return err
	}
	return s.base.Delete(id, ownerId)
}
