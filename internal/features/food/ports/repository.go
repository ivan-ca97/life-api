package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/food/domain"
)

type FoodRepository interface {
	Create(food *domain.Food) error
	FindById(id, userId uuid.UUID) (*domain.Food, error)
	List(userId uuid.UUID, params ListParams) (types.Page[domain.Food], error)
	Update(id, userId uuid.UUID, params UpdateParams) (*domain.Food, error)
	Delete(id, userId uuid.UUID) error
	Frequency(userId uuid.UUID, params FrequencyParams) ([]FrequencyResult, error)
	IngredientFrequency(userId uuid.UUID, params IngredientFrequencyParams) ([]IngredientFrequencyResult, error)
	ListIngredients(userId uuid.UUID, query *string) ([]domain.Ingredient, error)
}
