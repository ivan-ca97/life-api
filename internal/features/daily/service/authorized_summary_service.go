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

func (s *authorizedSummaryService) GetSummary(ctx context.Context, ownerId uuid.UUID, date time.Time) (*domain.DailySummary, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyRead)
	if err != nil {
		return nil, err
	}
	summary, err := s.base.GetSummary(ownerId, date)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (s *authorizedSummaryService) GetSummaryRange(ctx context.Context, ownerId uuid.UUID, from, to time.Time) ([]domain.DailySummary, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyRead)
	if err != nil {
		return nil, err
	}
	summaries, err := s.base.GetSummaryRange(ownerId, from, to)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}

func (s *authorizedSummaryService) GetDailyCheck(ctx context.Context, ownerId uuid.UUID, date time.Time) (*domain.DailyCheck, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyRead)
	if err != nil {
		return nil, err
	}
	check, err := s.base.GetDailyCheck(ownerId, date)
	if err != nil {
		return nil, err
	}
	return check, nil
}
