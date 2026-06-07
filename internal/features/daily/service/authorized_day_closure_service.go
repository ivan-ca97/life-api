package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/daily/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedDayClosureService struct {
	base       ports.DayClosureService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedDayClosureService = (*authorizedDayClosureService)(nil)

func NewAuthorizedDayClosureService(base ports.DayClosureService, authorizer auth.AuthorizationService) *authorizedDayClosureService {
	return &authorizedDayClosureService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedDayClosureService) Close(ctx context.Context, ownerId uuid.UUID, date time.Time) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate)
	if err != nil {
		return err
	}
	err = s.base.Close(ownerId, date)
	if err != nil {
		return err
	}
	return nil
}

func (s *authorizedDayClosureService) Open(ctx context.Context, ownerId uuid.UUID, date time.Time) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyUpdate)
	if err != nil {
		return err
	}
	err = s.base.Open(ownerId, date)
	if err != nil {
		return err
	}
	return nil
}

func (s *authorizedDayClosureService) IsClosed(ctx context.Context, ownerId uuid.UUID, date time.Time) (bool, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.DailyRead)
	if err != nil {
		return false, err
	}
	closed, err := s.base.IsClosed(ownerId, date)
	if err != nil {
		return false, err
	}
	return closed, nil
}
