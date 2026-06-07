package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/dayclosure"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
	"github.com/ivan-ca97/life/internal/features/exercise/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedExerciseService struct {
	base           ports.ExerciseService
	authorizer     auth.AuthorizationService
	closureChecker dayclosure.DayClosureChecker
}

var _ ports.AuthorizedExerciseService = (*authorizedExerciseService)(nil)

func NewAuthorizedExerciseService(base ports.ExerciseService, authorizer auth.AuthorizationService, closureChecker dayclosure.DayClosureChecker) *authorizedExerciseService {
	return &authorizedExerciseService{
		base:           base,
		authorizer:     authorizer,
		closureChecker: closureChecker,
	}
}

func (s *authorizedExerciseService) Create(ctx context.Context, ownerId uuid.UUID, params ports.CreateParams) (*domain.Exercise, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.ExercisesCreate)
	if err != nil {
		return nil, err
	}
	closed, err := s.closureChecker.IsClosed(ownerId, params.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}
	exercise, err := s.base.Create(ownerId, params)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *authorizedExerciseService) GetById(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) (*domain.Exercise, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.ExercisesRead)
	if err != nil {
		return nil, err
	}
	exercise, err := s.base.GetById(id, ownerId)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *authorizedExerciseService) List(ctx context.Context, ownerId uuid.UUID, params ports.ListParams) (types.Page[domain.Exercise], error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.ExercisesRead)
	if err != nil {
		return types.Page[domain.Exercise]{}, err
	}
	page, err := s.base.List(ownerId, params)
	if err != nil {
		return types.Page[domain.Exercise]{}, err
	}
	return page, nil
}

func (s *authorizedExerciseService) Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params ports.UpdateParams) (*domain.Exercise, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.ExercisesUpdate)
	if err != nil {
		return nil, err
	}
	exercise, err := s.base.GetById(id, ownerId)
	if err != nil {
		return nil, err
	}
	closed, err := s.closureChecker.IsClosed(ownerId, exercise.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}
	updated, err := s.base.Update(id, ownerId, params)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *authorizedExerciseService) Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.ExercisesDelete)
	if err != nil {
		return err
	}
	exercise, err := s.base.GetById(id, ownerId)
	if err != nil {
		return err
	}
	closed, err := s.closureChecker.IsClosed(ownerId, exercise.Date)
	if err != nil {
		return err
	}
	if closed {
		return dayclosure.ErrDayClosed
	}
	err = s.base.Delete(id, ownerId)
	if err != nil {
		return err
	}
	return nil
}
