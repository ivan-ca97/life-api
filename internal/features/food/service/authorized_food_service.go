package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/food/domain"
	"github.com/ivan-ca97/life/internal/features/food/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedFoodService struct {
	base       ports.FoodService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedFoodService = (*authorizedFoodService)(nil)

func NewAuthorizedFoodService(base ports.FoodService, authorizer auth.AuthorizationService) *authorizedFoodService {
	return &authorizedFoodService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedFoodService) Create(ctx context.Context, params ports.CreateParams) (*domain.Food, error) {
	err := s.authorizer.Require(ctx, permissions.FoodsCreate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	food, err := s.base.Create(userId, params)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *authorizedFoodService) GetById(ctx context.Context, id uuid.UUID) (*domain.Food, error) {
	err := s.authorizer.Require(ctx, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	food, err := s.base.GetById(id, userId)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *authorizedFoodService) List(ctx context.Context, params ports.ListParams) (types.Page[domain.Food], error) {
	err := s.authorizer.Require(ctx, permissions.FoodsRead)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	page, err := s.base.List(userId, params)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	return page, nil
}

func (s *authorizedFoodService) Update(ctx context.Context, id uuid.UUID, params ports.UpdateParams) (*domain.Food, error) {
	err := s.authorizer.Require(ctx, permissions.FoodsUpdate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	food, err := s.base.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *authorizedFoodService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.authorizer.Require(ctx, permissions.FoodsDelete)
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

func (s *authorizedFoodService) Frequency(ctx context.Context, params ports.FrequencyParams) ([]ports.FrequencyResult, error) {
	err := s.authorizer.Require(ctx, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	results, err := s.base.Frequency(userId, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *authorizedFoodService) IngredientFrequency(ctx context.Context, params ports.IngredientFrequencyParams) ([]ports.IngredientFrequencyResult, error) {
	err := s.authorizer.Require(ctx, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	results, err := s.base.IngredientFrequency(userId, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *authorizedFoodService) ListIngredients(ctx context.Context, query *string) ([]domain.Ingredient, error) {
	err := s.authorizer.Require(ctx, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.base.ListIngredients(userId, query)
}
