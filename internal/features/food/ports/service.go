package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/food/domain"
)

type ConversionParam struct {
	Unit           string
	BaseEquivalent float64
	Inverse        bool
	Note           *string
}

type CreateParams struct {
	Name                string
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	MeasurementType     string
	BaseQuantity        float64
	BaseUnit            string
	Public              bool
	Tags                []string
	Ingredients         []string
	Conversions         []ConversionParam
}

type UpdateParams struct {
	Name                *string
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	MeasurementType     *string
	BaseQuantity        *float64
	BaseUnit            *string
	Public              *bool
	Tags                *[]string
	Ingredients         *[]string
	Conversions         *[]ConversionParam
}

type ListParams struct {
	types.PaginationParams
	Query *string
	Tag   *string
}

type FrequencyParams struct {
	From *time.Time
	To   *time.Time
	Tag  *string
}

type FrequencyResult struct {
	FoodId   uuid.UUID
	FoodName string
	Count    int64
	LastDate time.Time
}

type IngredientFrequencyParams struct {
	From *time.Time
	To   *time.Time
}

type IngredientFrequencyResult struct {
	IngredientId   uuid.UUID
	IngredientName string
	Count          int64
	LastDate       time.Time
}

type CommunityListParams struct {
	types.PaginationParams
	Query *string
}

type FoodService interface {
	Create(userId uuid.UUID, params CreateParams) (*domain.Food, error)
	GetById(id, userId uuid.UUID) (*domain.Food, error)
	List(userId uuid.UUID, params ListParams) (types.Page[domain.Food], error)
	Update(id, userId uuid.UUID, params UpdateParams) (*domain.Food, error)
	Delete(id, userId uuid.UUID) error
	Frequency(userId uuid.UUID, params FrequencyParams) ([]FrequencyResult, error)
	IngredientFrequency(userId uuid.UUID, params IngredientFrequencyParams) ([]IngredientFrequencyResult, error)
	ListIngredients(userId uuid.UUID, query *string) ([]domain.Ingredient, error)
	ListCommunity(params CommunityListParams) (types.Page[domain.Food], error)
	Copy(actorId, foodId uuid.UUID) (*domain.Food, error)
}
