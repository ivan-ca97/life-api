package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
)

type ItemParam struct {
	FoodId            uuid.UUID
	Quantity          float64
	Unit              string
	Notes             string
	MeasurementMethod domain.MeasurementMethod
}

type PhotoParam struct {
	Url        string
	IsPrimary  bool
	MealItemId *uuid.UUID
}

type CreateParams struct {
	Date         time.Time
	Type         string
	Name         string
	Photos       []PhotoParam
	EatenAt      *time.Time
	Calories     *float64
	ProteinGrams *float64
	CarbsGrams   *float64
	FatGrams     *float64
	FiberGrams   *float64
	Tags         []string
	Items        []ItemParam
	Notes        string
}

type UpdateParams struct {
	Date          *time.Time
	Type          *string
	Name          *string
	Photos        *[]PhotoParam
	EatenAt       *time.Time
	Calories      *float64
	ProteinGrams  *float64
	CarbsGrams    *float64
	FatGrams      *float64
	FiberGrams    *float64
	Tags          *[]string
	Items         *[]ItemParam
	ResolvedItems *[]domain.MealItem
	Notes         *string
}

type ListParams struct {
	types.PaginationParams
	Date *time.Time
}

type NutritionPreviewItem struct {
	FoodId       uuid.UUID
	Calories     *float64
	ProteinGrams *float64
	CarbsGrams   *float64
	FatGrams     *float64
	FiberGrams   *float64
}

type NutritionPreview struct {
	Calories     *float64
	ProteinGrams *float64
	CarbsGrams   *float64
	FatGrams     *float64
	FiberGrams   *float64
	Items        []NutritionPreviewItem
}

type MealService interface {
	Create(userId uuid.UUID, params CreateParams) (*domain.Meal, error)
	GetById(id, userId uuid.UUID) (*domain.Meal, error)
	List(userId uuid.UUID, params ListParams) (types.Page[domain.Meal], error)
	Update(id, userId uuid.UUID, params UpdateParams) (*domain.Meal, error)
	Delete(id, userId uuid.UUID) error
	ListTypes(userId uuid.UUID, hour *int) ([]string, error)
	PreviewNutrition(userId uuid.UUID, items []ItemParam) (*NutritionPreview, error)
}
