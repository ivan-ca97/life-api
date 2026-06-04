package ports

import "github.com/google/uuid"

type FoodConversion struct {
	Unit           string
	BaseEquivalent float64
	Inverse        bool
}

type FoodNutrition struct {
	Id                  uuid.UUID
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	MeasurementType     string
	BaseQuantity        float64
	BaseUnit            string
	Conversions         []FoodConversion
}

type FoodLookup interface {
	FindByIds(userId uuid.UUID, ids []uuid.UUID) (map[uuid.UUID]FoodNutrition, error)
}
