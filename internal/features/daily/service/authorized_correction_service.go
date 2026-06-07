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

type authorizedCorrectionService struct {
	base       ports.CorrectionService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedCorrectionService = (*authorizedCorrectionService)(nil)

func NewAuthorizedCorrectionService(base ports.CorrectionService, authorizer auth.AuthorizationService) *authorizedCorrectionService {
	return &authorizedCorrectionService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedCorrectionService) GetCorrection(ctx context.Context, ownerId uuid.UUID, date time.Time) (*domain.Correction, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyRead)
	if err != nil {
		return nil, err
	}
	correction, err := s.base.GetCorrection(ownerId, date)
	if err != nil {
		return nil, err
	}
	return correction, nil
}

func (s *authorizedCorrectionService) UpsertCorrection(ctx context.Context, ownerId uuid.UUID, correction *domain.Correction) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate)
	if err != nil {
		return err
	}
	err = s.base.UpsertCorrection(ownerId, correction)
	if err != nil {
		return err
	}
	return nil
}

func (s *authorizedCorrectionService) DeleteCorrection(ctx context.Context, ownerId uuid.UUID, date time.Time) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate)
	if err != nil {
		return err
	}
	err = s.base.DeleteCorrection(ownerId, date)
	if err != nil {
		return err
	}
	return nil
}
