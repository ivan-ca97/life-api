package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
)

type VolumeConversionParam struct {
	GramsPerMl float64
	Note       *string
}

type ConversionsParam struct {
	VolumeConversion *VolumeConversionParam
}

type PortionParam struct {
	Name           string
	BaseEquivalent float64
}

type UpdatePortionParams struct {
	Name           *string
	BaseEquivalent *float64
}

type CreateFoodParams struct {
	Name                string
	PhotoUrl            string
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	BaseQuantity        float64
	BaseUnit            string
	Tags                []string
	Ingredients         []string
	Conversions         *ConversionsParam
	Portions            []PortionParam
}

type UpdateFoodParams struct {
	Name                *string
	PhotoUrl            *string
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	BaseQuantity        *float64
	BaseUnit            *string
	Tags                *[]string
	Ingredients         *[]string
	Conversions         *ConversionsParam
}

type ListFoodParams struct {
	types.PaginationParams
	Query *string
	Tag   *string
	Sort  *string // "name" (default), "created_at", "updated_at"
}

type CommunityListParams struct {
	types.PaginationParams
	Query *string
}

type FoodService interface {
	Create(userId uuid.UUID, params CreateFoodParams) (*domain.Food, error)
	GetById(id, userId uuid.UUID) (*domain.Food, error)
	FindUsableByIds(userId uuid.UUID, ids []uuid.UUID) ([]domain.Food, error)
	List(userId uuid.UUID, params ListFoodParams) (types.Page[domain.Food], error)
	Update(id, userId uuid.UUID, params UpdateFoodParams) (*domain.Food, error)
	Delete(id, userId uuid.UUID) error
	AddPortion(foodId, userId uuid.UUID, params PortionParam) (*domain.Portion, error)
	UpdatePortion(foodId, portionId, userId uuid.UUID, params UpdatePortionParams) (*domain.Portion, error)
	DeletePortion(foodId, portionId, userId uuid.UUID) error
	ListIngredients(userId uuid.UUID, query *string) ([]domain.Ingredient, error)
	ListCommunity(params CommunityListParams) (types.Page[domain.Food], error)
	Copy(actorId, foodId uuid.UUID) (*domain.Food, error)
}
