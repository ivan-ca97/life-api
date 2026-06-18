package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/food/domain"
)

type AuthorizedFoodService interface {
	Create(ctx context.Context, ownerId uuid.UUID, params CreateParams) (*domain.Food, error)
	GetById(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) (*domain.Food, error)
	List(ctx context.Context, ownerId uuid.UUID, params ListParams) (types.Page[domain.Food], error)
	Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params UpdateParams) (*domain.Food, error)
	Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error
	Frequency(ctx context.Context, ownerId uuid.UUID, params FrequencyParams) ([]FrequencyResult, error)
	IngredientFrequency(ctx context.Context, ownerId uuid.UUID, params IngredientFrequencyParams) ([]IngredientFrequencyResult, error)
	ListIngredients(ctx context.Context, ownerId uuid.UUID, query *string) ([]domain.Ingredient, error)
	ListCommunity(ctx context.Context, params CommunityListParams) (types.Page[domain.Food], error)
	Copy(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID) (*domain.Food, error)
	Impact(ctx context.Context, ownerId uuid.UUID, foodId uuid.UUID) (*ImpactResult, error)
}
