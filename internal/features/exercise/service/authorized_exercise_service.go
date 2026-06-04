package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
	"github.com/ivan-ca97/life/internal/features/exercise/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedExerciseService struct {
	base       ports.ExerciseService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedExerciseService = (*authorizedExerciseService)(nil)

func NewAuthorizedExerciseService(base ports.ExerciseService, authorizer auth.AuthorizationService) *authorizedExerciseService {
	return &authorizedExerciseService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedExerciseService) Create(ctx context.Context, params ports.CreateParams) (*domain.Exercise, error) {
	err := s.authorizer.Require(ctx, permissions.ExercisesCreate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	exercise, err := s.base.Create(userId, params)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *authorizedExerciseService) GetById(ctx context.Context, id uuid.UUID) (*domain.Exercise, error) {
	err := s.authorizer.Require(ctx, permissions.ExercisesRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	exercise, err := s.base.GetById(id, userId)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *authorizedExerciseService) List(ctx context.Context, params ports.ListParams) (types.Page[domain.Exercise], error) {
	err := s.authorizer.Require(ctx, permissions.ExercisesRead)
	if err != nil {
		return types.Page[domain.Exercise]{}, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return types.Page[domain.Exercise]{}, err
	}
	page, err := s.base.List(userId, params)
	if err != nil {
		return types.Page[domain.Exercise]{}, err
	}
	return page, nil
}

func (s *authorizedExerciseService) Update(ctx context.Context, id uuid.UUID, params ports.UpdateParams) (*domain.Exercise, error) {
	err := s.authorizer.Require(ctx, permissions.ExercisesUpdate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	exercise, err := s.base.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *authorizedExerciseService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.authorizer.Require(ctx, permissions.ExercisesDelete)
	if err != nil {
		return err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return err
	}
	err = s.base.Delete(id, userId)
	if err != nil {
		return err
	}
	return nil
}
