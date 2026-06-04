package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
	"github.com/ivan-ca97/life/internal/features/meal/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedMealService struct {
	base       ports.MealService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedMealService = (*authorizedMealService)(nil)

func NewAuthorizedMealService(base ports.MealService, authorizer auth.AuthorizationService) *authorizedMealService {
	return &authorizedMealService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedMealService) Create(ctx context.Context, params ports.CreateParams) (*domain.Meal, error) {
	err := s.authorizer.Require(ctx, permissions.MealsCreate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	meal, err := s.base.Create(userId, params)
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (s *authorizedMealService) GetById(ctx context.Context, id uuid.UUID) (*domain.Meal, error) {
	err := s.authorizer.Require(ctx, permissions.MealsRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	meal, err := s.base.GetById(id, userId)
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (s *authorizedMealService) List(ctx context.Context, params ports.ListParams) (types.Page[domain.Meal], error) {
	err := s.authorizer.Require(ctx, permissions.MealsRead)
	if err != nil {
		return types.Page[domain.Meal]{}, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return types.Page[domain.Meal]{}, err
	}
	page, err := s.base.List(userId, params)
	if err != nil {
		return types.Page[domain.Meal]{}, err
	}
	return page, nil
}

func (s *authorizedMealService) Update(ctx context.Context, id uuid.UUID, params ports.UpdateParams) (*domain.Meal, error) {
	err := s.authorizer.Require(ctx, permissions.MealsUpdate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	meal, err := s.base.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (s *authorizedMealService) PreviewNutrition(ctx context.Context, items []ports.ItemParam) (*ports.NutritionPreview, error) {
	err := s.authorizer.Require(ctx, permissions.MealsRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.base.PreviewNutrition(userId, items)
}

func (s *authorizedMealService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.authorizer.Require(ctx, permissions.MealsDelete)
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

func (s *authorizedMealService) ListTypes(ctx context.Context, hour *int) ([]string, error) {
	err := s.authorizer.Require(ctx, permissions.MealsRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	mealTypes, err := s.base.ListTypes(userId, hour)
	if err != nil {
		return nil, err
	}
	return mealTypes, nil
}
