package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
)

type FoodRepository interface {
	Create(food *domain.Food) error
	FindById(id, userId uuid.UUID) (*domain.Food, error)
	FindByIdGlobal(id uuid.UUID) (*domain.Food, error)
	FindUsableByIds(userId uuid.UUID, ids []uuid.UUID) ([]domain.Food, error)
	List(userId uuid.UUID, params ListFoodParams) (types.Page[domain.Food], error)
	ListCommunity(params CommunityListParams) (types.Page[domain.Food], error)
	Update(id, userId uuid.UUID, params UpdateFoodParams) (*domain.Food, error)
	Delete(id, userId uuid.UUID) error
	AddPortion(foodId, userId uuid.UUID, portion *domain.Portion) error
	UpdatePortion(foodId, userId uuid.UUID, portion *domain.Portion) error
	DeletePortion(foodId, portionId, userId uuid.UUID) error
	Frequency(userId uuid.UUID, params FrequencyParams) ([]FrequencyResult, error)
	IngredientFrequency(userId uuid.UUID, params IngredientFrequencyParams) ([]IngredientFrequencyResult, error)
	ListIngredients(userId uuid.UUID, query *string) ([]domain.Ingredient, error)
	IsAccessibleBy(foodId, actorId uuid.UUID) (bool, error)
	Impact(foodId uuid.UUID) (*ImpactResult, error)
}
