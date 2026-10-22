package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
	"github.com/ivan-ca97/life/internal/features/nutrition/ports"
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

func (s *authorizedMealService) Create(ctx context.Context, ownerId uuid.UUID, params ports.CreateMealParams) (*domain.Meal, error) {
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

func (s *authorizedMealService) List(ctx context.Context, ownerId uuid.UUID, params ports.ListMealParams) (types.Page[domain.Meal], error) {
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

func (s *authorizedMealService) Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params ports.UpdateMealParams) (*domain.Meal, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MealsUpdate)
	if err != nil {
		return nil, err
	}
	updated, err := s.base.Update(id, ownerId, params)
	if err != nil {
		return nil, err
	}
	return updated, nil
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

func (s *authorizedMealService) FoodFrequency(ctx context.Context, ownerId uuid.UUID, params ports.FoodFrequencyParams) ([]ports.FoodFrequencyResult, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	results, err := s.base.FoodFrequency(ownerId, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *authorizedMealService) IngredientFrequency(ctx context.Context, ownerId uuid.UUID, params ports.IngredientFrequencyParams) ([]ports.IngredientFrequencyResult, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	results, err := s.base.IngredientFrequency(ownerId, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *authorizedMealService) FoodImpact(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID) (*ports.FoodImpactResult, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	impact, err := s.base.FoodImpact(foodId, ownerId)
	if err != nil {
		return nil, err
	}
	return impact, nil
}
