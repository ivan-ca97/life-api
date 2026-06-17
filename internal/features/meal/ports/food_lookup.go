package ports

import "github.com/google/uuid"

type FoodVolumeConversion struct {
	GramsPerMl float64
}

type FoodUnitConversion struct {
	BaseEquivalent float64
}

type FoodPortion struct {
	Id             uuid.UUID
	Name           string
	BaseEquivalent float64
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
	VolumeConversion    *FoodVolumeConversion
	UnitConversion      *FoodUnitConversion
	Portions            []FoodPortion
}

type FoodLookup interface {
	FindByIds(userId uuid.UUID, ids []uuid.UUID) (map[uuid.UUID]FoodNutrition, error)
}
