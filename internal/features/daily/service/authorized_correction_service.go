package service

import (
	"context"
	"time"

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

func (s *authorizedCorrectionService) GetCorrection(ctx context.Context, date time.Time) (*domain.Correction, error) {
	if err := s.authorizer.Require(ctx, permissions.DailyRead); err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.base.GetCorrection(userId, date)
}

func (s *authorizedCorrectionService) UpsertCorrection(ctx context.Context, correction *domain.Correction) error {
	if err := s.authorizer.Require(ctx, permissions.DailyUpdate); err != nil {
		return err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return err
	}
	return s.base.UpsertCorrection(userId, correction)
}

func (s *authorizedCorrectionService) DeleteCorrection(ctx context.Context, date time.Time) error {
	if err := s.authorizer.Require(ctx, permissions.DailyUpdate); err != nil {
		return err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return err
	}
	return s.base.DeleteCorrection(userId, date)
}
