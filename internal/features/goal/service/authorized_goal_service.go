package service

import (
	"context"

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

func (s *authorizedGoalService) GetCurrent(ctx context.Context) (*domain.Goal, error) {
	err := s.authorizer.Require(ctx, permissions.GoalsRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	goal, err := s.base.GetByUserId(userId)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *authorizedGoalService) Upsert(ctx context.Context, params ports.UpsertParams) (*domain.Goal, error) {
	err := s.authorizer.Require(ctx, permissions.GoalsUpdate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	goal, err := s.base.Upsert(userId, params)
	if err != nil {
		return nil, err
	}
	return goal, nil
}
