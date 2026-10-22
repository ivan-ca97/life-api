package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
)

type AuthorizedFoodService interface {
	Create(ctx context.Context, ownerId uuid.UUID, params CreateFoodParams) (*domain.Food, error)
	GetById(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) (*domain.Food, error)
	List(ctx context.Context, ownerId uuid.UUID, params ListFoodParams) (types.Page[domain.Food], error)
	Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params UpdateFoodParams) (*domain.Food, error)
	Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error
	AddPortion(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID, params PortionParam) (*domain.Portion, error)
	UpdatePortion(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID, portionId uuid.UUID, params UpdatePortionParams) (*domain.Portion, error)
	DeletePortion(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID, portionId uuid.UUID) error
	ListIngredients(ctx context.Context, ownerId uuid.UUID, query *string) ([]domain.Ingredient, error)
	ListCommunity(ctx context.Context, params CommunityListParams) (types.Page[domain.Food], error)
	Copy(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID) (*domain.Food, error)
}
