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

func (s *authorizedMealService) Create(ctx context.Context, ownerId uuid.UUID, params ports.CreateParams) (*domain.Meal, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MealsCreate)
	if err != nil {
		return nil, err
	}
	meal, err := s.base.Create(ownerId, params)
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (s *authorizedMealService) GetById(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) (*domain.Meal, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MealsRead)
	if err != nil {
		return nil, err
	}
	meal, err := s.base.GetById(id, ownerId)
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (s *authorizedMealService) List(ctx context.Context, ownerId uuid.UUID, params ports.ListParams) (types.Page[domain.Meal], error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MealsRead)
	if err != nil {
		return types.Page[domain.Meal]{}, err
	}
	page, err := s.base.List(ownerId, params)
	if err != nil {
		return types.Page[domain.Meal]{}, err
	}
	return page, nil
}

func (s *authorizedMealService) Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params ports.UpdateParams) (*domain.Meal, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MealsUpdate)
	if err != nil {
		return nil, err
	}
	meal, err := s.base.Update(id, ownerId, params)
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (s *authorizedMealService) Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MealsDelete)
	if err != nil {
		return err
	}
	err = s.base.Delete(id, ownerId)
	if err != nil {
		return err
	}
	return nil
}

func (s *authorizedMealService) ListTypes(ctx context.Context, ownerId uuid.UUID, hour *int) ([]string, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MealsRead)
	if err != nil {
		return nil, err
	}
	types, err := s.base.ListTypes(ownerId, hour)
	if err != nil {
		return nil, err
	}
	return types, nil
}

func (s *authorizedMealService) PreviewNutrition(ctx context.Context, ownerId uuid.UUID, items []ports.ItemParam) (*ports.NutritionPreview, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MealsRead)
	if err != nil {
		return nil, err
	}
	preview, err := s.base.PreviewNutrition(ownerId, items)
	if err != nil {
		return nil, err
	}
	return preview, nil
}
