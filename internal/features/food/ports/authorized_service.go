package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/food/domain"
)

type AuthorizedFoodService interface {
	Create(ctx context.Context, params CreateParams) (*domain.Food, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Food, error)
	List(ctx context.Context, params ListParams) (types.Page[domain.Food], error)
	Update(ctx context.Context, id uuid.UUID, params UpdateParams) (*domain.Food, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Frequency(ctx context.Context, params FrequencyParams) ([]FrequencyResult, error)
	IngredientFrequency(ctx context.Context, params IngredientFrequencyParams) ([]IngredientFrequencyResult, error)
	ListIngredients(ctx context.Context, query *string) ([]domain.Ingredient, error)
}
