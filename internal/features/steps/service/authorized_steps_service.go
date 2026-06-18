package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/steps/domain"
	"github.com/ivan-ca97/life/internal/features/steps/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedStepsService struct {
	base       ports.StepsService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedStepsService = (*authorizedStepsService)(nil)

func NewAuthorizedStepsService(base ports.StepsService, authorizer auth.AuthorizationService) *authorizedStepsService {
	return &authorizedStepsService{base: base, authorizer: authorizer}
}

func (s *authorizedStepsService) Upsert(ctx context.Context, userId uuid.UUID, date time.Time, params ports.UpsertParams) (*domain.DailySteps, error) {
	if err := s.authorizer.Authorize(ctx, userId, permissions.StepsWrite); err != nil {
		return nil, err
	}
	return s.base.Upsert(userId, date, params)
}

func (s *authorizedStepsService) GetByDate(ctx context.Context, userId uuid.UUID, date time.Time) (*domain.DailySteps, error) {
	if err := s.authorizer.Authorize(ctx, userId, permissions.StepsRead); err != nil {
		return nil, err
	}
	return s.base.GetByDate(userId, date)
}

func (s *authorizedStepsService) List(ctx context.Context, userId uuid.UUID, params ports.ListParams) ([]domain.DailySteps, error) {
	if err := s.authorizer.Authorize(ctx, userId, permissions.StepsRead); err != nil {
		return nil, err
	}
	return s.base.List(userId, params)
}

func (s *authorizedStepsService) Delete(ctx context.Context, userId uuid.UUID, date time.Time) error {
	if err := s.authorizer.Authorize(ctx, userId, permissions.StepsWrite); err != nil {
		return err
	}
	return s.base.Delete(userId, date)
}
