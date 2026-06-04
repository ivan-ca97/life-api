package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
)

type AuthorizedMealService interface {
	Create(ctx context.Context, params CreateParams) (*domain.Meal, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Meal, error)
	List(ctx context.Context, params ListParams) (types.Page[domain.Meal], error)
	Update(ctx context.Context, id uuid.UUID, params UpdateParams) (*domain.Meal, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListTypes(ctx context.Context, hour *int) ([]string, error)
	PreviewNutrition(ctx context.Context, items []ItemParam) (*NutritionPreview, error)
}
