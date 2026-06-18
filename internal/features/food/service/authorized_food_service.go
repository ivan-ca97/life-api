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

func (s *authorizedFoodService) Create(ctx context.Context, ownerId uuid.UUID, params ports.CreateParams) (*domain.Food, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsCreate)
	if err != nil {
		return nil, err
	}
	food, err := s.base.Create(ownerId, params)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *authorizedFoodService) GetById(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) (*domain.Food, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	food, err := s.base.GetById(id, ownerId)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *authorizedFoodService) List(ctx context.Context, ownerId uuid.UUID, params ports.ListParams) (types.Page[domain.Food], error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsRead)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	page, err := s.base.List(ownerId, params)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	return page, nil
}

func (s *authorizedFoodService) Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params ports.UpdateParams) (*domain.Food, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsUpdate)
	if err != nil {
		return nil, err
	}
	food, err := s.base.Update(id, ownerId, params)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *authorizedFoodService) Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsDelete)
	if err != nil {
		return err
	}
	err = s.base.Delete(id, ownerId)
	if err != nil {
		return err
	}
	return nil
}

func (s *authorizedFoodService) Frequency(ctx context.Context, ownerId uuid.UUID, params ports.FrequencyParams) ([]ports.FrequencyResult, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	results, err := s.base.Frequency(ownerId, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *authorizedFoodService) IngredientFrequency(ctx context.Context, ownerId uuid.UUID, params ports.IngredientFrequencyParams) ([]ports.IngredientFrequencyResult, error) {
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

func (s *authorizedFoodService) ListIngredients(ctx context.Context, ownerId uuid.UUID, query *string) ([]domain.Ingredient, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	ingredients, err := s.base.ListIngredients(ownerId, query)
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (s *authorizedFoodService) ListCommunity(ctx context.Context, params ports.CommunityListParams) (types.Page[domain.Food], error) {
	_, err := auth.ActorFromContext(ctx)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	page, err := s.base.ListCommunity(params)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	return page, nil
}

func (s *authorizedFoodService) Copy(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID) (*domain.Food, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsCreate)
	if err != nil {
		return nil, err
	}
	food, err := s.base.Copy(ownerId, foodId)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *authorizedFoodService) Impact(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID) (*ports.ImpactResult, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.FoodsRead)
	if err != nil {
		return nil, err
	}
	return s.base.Impact(foodId, ownerId)
}
