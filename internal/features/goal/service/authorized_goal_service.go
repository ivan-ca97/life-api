package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
	"github.com/ivan-ca97/life/internal/features/goal/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedGoalService struct {
	base       ports.GoalService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedGoalService = (*authorizedGoalService)(nil)

func NewAuthorizedGoalService(base ports.GoalService, authorizer auth.AuthorizationService) *authorizedGoalService {
	return &authorizedGoalService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedGoalService) GetCurrent(ctx context.Context, ownerId uuid.UUID) (*domain.Goal, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.GoalsRead)
	if err != nil {
		return nil, err
	}
	goal, err := s.base.GetByUserId(ownerId)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *authorizedGoalService) Upsert(ctx context.Context, ownerId uuid.UUID, params ports.UpsertParams) (*domain.Goal, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.GoalsUpdate)
	if err != nil {
		return nil, err
	}
	goal, err := s.base.Upsert(ownerId, params)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *authorizedGoalService) GetProgress(ctx context.Context, ownerId uuid.UUID, from, to time.Time) (*domain.GoalProgress, error) {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.GoalsRead); err != nil {
		return nil, err
	}
	return s.base.GetProgress(ownerId, from, to)
}
