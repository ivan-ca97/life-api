package domain

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/units"
)

// Builds the item that results from eating quantity of food.
func NewMealItem(food Food, quantity float64, unit string) (*MealItem, error) {
	unitDimension, ok := units.GetDimension(unit)
	if !ok {
		return nil, NewUnknownUnitError(food.Id, unit)
	}
	inStandardUnit, _, _ := units.ConvertToStandard(quantity, unit)

	var normalized float64
	switch {
	case unitDimension == food.Dimension():
		normalized = inStandardUnit
	case canConvertByDensity(food, unitDimension):
		normalized = convertByDensity(food, inStandardUnit)
	default:
		return nil, NewUnknownUnitError(food.Id, unit)
	}
	item := newMealItem(food, quantity, unit, normalized)
	return item, nil
}

// Builds the item that results from eating count of one of the food's portions.
// The item keeps the portion's name and equivalent, so editing the portion later
// does not change what was eaten.
func NewMealItemOfPortion(food Food, count float64, portionId uuid.UUID) (*MealItem, error) {
	for _, portion := range food.Portions {
		if portion.Id != portionId {
			continue
		}

		item := newMealItem(food, count, portion.Name, count*portion.BaseEquivalent)
		id := portion.Id
		equivalent := portion.BaseEquivalent
		item.InputPortionId = &id
		item.InputPortionEquivalent = &equivalent
		return item, nil
	}
	return nil, NewUnknownPortionError(food.Id, portionId)
}

func newMealItem(food Food, inputQuantity float64, inputUnit string, normalizedQuantity float64) *MealItem {
	baseInStandard, _, _ := units.ConvertToStandard(food.BaseQuantity, food.BaseUnit)
	ratio := normalizedQuantity / baseInStandard

	item := &MealItem{
		FoodId:             food.Id,
		InputQuantity:      inputQuantity,
		InputUnit:          inputUnit,
		NormalizedQuantity: normalizedQuantity,
		NormalizedUnit:     units.StandardUnit[food.Dimension()],
		Calories:           scaleNutrient(food.DefaultCalories, ratio),
		ProteinGrams:       scaleNutrient(food.DefaultProteinGrams, ratio),
		CarbsGrams:         scaleNutrient(food.DefaultCarbsGrams, ratio),
		FatGrams:           scaleNutrient(food.DefaultFatGrams, ratio),
		FiberGrams:         scaleNutrient(food.DefaultFiberGrams, ratio),
	}
	return item
}

func scaleNutrient(amount *float64, ratio float64) *float64 {
	if amount == nil {
		return nil
	}
	scaled := *amount * ratio
	return &scaled
}

// Mass and volume convert into each other when the food has a density.
func canConvertByDensity(food Food, unitDimension units.Dimension) bool {
	if food.VolumeConversion == nil {
		return false
	}
	dimension := food.Dimension()
	massToVolume := dimension == units.DimensionVolume && unitDimension == units.DimensionMass
	volumeToMass := dimension == units.DimensionMass && unitDimension == units.DimensionVolume
	return massToVolume || volumeToMass
}

// Takes a quantity of the other dimension, in its standard unit, to the food's
// standard unit.
func convertByDensity(food Food, quantity float64) float64 {
	gramsPerMl := food.VolumeConversion.GramsPerMl
	if food.Dimension() == units.DimensionMass {
		return quantity * gramsPerMl
	}
	return quantity / gramsPerMl
}
