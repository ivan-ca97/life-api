package service

import (
	"context"
	"time"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedSummaryService struct {
	base       ports.SummaryService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedSummaryService = (*authorizedSummaryService)(nil)

func NewAuthorizedSummaryService(base ports.SummaryService, authorizer auth.AuthorizationService) *authorizedSummaryService {
	return &authorizedSummaryService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedSummaryService) GetSummary(ctx context.Context, date time.Time) (*domain.DailySummary, error) {
	err := s.authorizer.Require(ctx, permissions.DailyRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	summary, err := s.base.GetSummary(userId, date)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (s *authorizedSummaryService) GetSummaryRange(ctx context.Context, from, to time.Time) ([]domain.DailySummary, error) {
	err := s.authorizer.Require(ctx, permissions.DailyRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	summaries, err := s.base.GetSummaryRange(userId, from, to)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}
