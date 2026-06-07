package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
)

type AuthorizedMealService interface {
	Create(ctx context.Context, ownerId uuid.UUID, params CreateParams) (*domain.Meal, error)
	GetById(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) (*domain.Meal, error)
	List(ctx context.Context, ownerId uuid.UUID, params ListParams) (types.Page[domain.Meal], error)
	Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params UpdateParams) (*domain.Meal, error)
	Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error
	ListTypes(ctx context.Context, ownerId uuid.UUID, hour *int) ([]string, error)
	PreviewNutrition(ctx context.Context, ownerId uuid.UUID, items []ItemParam) (*NutritionPreview, error)
}
