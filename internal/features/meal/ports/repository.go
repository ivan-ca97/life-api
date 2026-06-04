package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
)

type MealRepository interface {
	Create(meal *domain.Meal) error
	FindById(id, userId uuid.UUID) (*domain.Meal, error)
	List(userId uuid.UUID, params ListParams) (types.Page[domain.Meal], error)
	Update(id, userId uuid.UUID, params UpdateParams) (*domain.Meal, error)
	Delete(id, userId uuid.UUID) error
	ListDistinctTypes(userId uuid.UUID, hour *int) ([]string, error)
}
