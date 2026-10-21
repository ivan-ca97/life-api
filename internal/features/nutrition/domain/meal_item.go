package domain

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/units"
)

// Quantity expressed in the food's standard unit, together
// with the named portion it came from, if any.
type resolvedUnit struct {
	NormalizedQuantity float64
	NormalizedUnit     string
	PortionId          *uuid.UUID
	PortionEquivalent  *float64
}

// Builds the item that results from eating quantity of food.
// Notes and measurement method are left empty for the caller to fill in.
func NewMealItem(food Food, quantity float64, unit string) (*MealItem, error) {
	resolved, err := resolveUnit(food, quantity, unit)
	if err != nil {
		return nil, err
	}

	baseInStandard, _, _ := units.ConvertToStandard(food.BaseQuantity, food.BaseUnit)
	ratio := resolved.NormalizedQuantity / baseInStandard

	item := &MealItem{
		FoodId:                 food.Id,
		InputQuantity:          quantity,
		InputUnit:              unit,
		InputPortionId:         resolved.PortionId,
		InputPortionEquivalent: resolved.PortionEquivalent,
		NormalizedQuantity:     resolved.NormalizedQuantity,
		NormalizedUnit:         resolved.NormalizedUnit,
		Calories:               scaleNutrient(food.DefaultCalories, ratio),
		ProteinGrams:           scaleNutrient(food.DefaultProteinGrams, ratio),
		CarbsGrams:             scaleNutrient(food.DefaultCarbsGrams, ratio),
		FatGrams:               scaleNutrient(food.DefaultFatGrams, ratio),
		FiberGrams:             scaleNutrient(food.DefaultFiberGrams, ratio),
	}
	return item, nil
}

// An undefined amount stays undefined.
func scaleNutrient(amount *float64, ratio float64) *float64 {
	if amount == nil {
		return nil
	}
	scaled := *amount * ratio
	return &scaled
}

// Reads the unit in the order below and fails when none of them applies: a metric
// unit of the food's dimension, a unit of the other dimension through the food's
// density, or one of the food's named portions.
func resolveUnit(food Food, quantity float64, unit string) (*resolvedUnit, error) {
	resolved := resolveMetricUnit(food, quantity, unit)
	if resolved != nil {
		return resolved, nil
	}
	resolved = resolveByDensity(food, quantity, unit)
	if resolved != nil {
		return resolved, nil
	}
	resolved = resolveNamedPortion(food, quantity, unit)
	if resolved != nil {
		return resolved, nil
	}
	return nil, NewUnknownUnitError(food.Id, unit)
}

// Reads the unit as a metric unit of the food's own dimension, such as grams for a
// mass food. Returns nil when it is not one.
func resolveMetricUnit(food Food, quantity float64, unit string) *resolvedUnit {
	dimension := units.Dimension(food.MeasurementType)
	unitDimension, ok := units.GetDimension(unit)
	if !ok || unitDimension != dimension {
		return nil
	}

	converted, _, _ := units.ConvertToStandard(quantity, unit)
	resolved := &resolvedUnit{
		NormalizedQuantity: converted,
		NormalizedUnit:     units.StandardUnit[dimension],
	}
	return resolved
}

// Reads the unit as one of the other dimension, going through the food's density:
// a volume for a mass food, a mass for a volume food. Returns nil when the food has
// no density or the unit is not of the other dimension.
func resolveByDensity(food Food, quantity float64, unit string) *resolvedUnit {
	if food.VolumeConversion == nil {
		return nil
	}
	unitDimension, ok := units.GetDimension(unit)
	if !ok {
		return nil
	}

	dimension := units.Dimension(food.MeasurementType)
	converted, _, _ := units.ConvertToStandard(quantity, unit)
	gramsPerMl := food.VolumeConversion.GramsPerMl

	var normalized float64
	switch {
	case dimension == units.DimensionMass && unitDimension == units.DimensionVolume:
		normalized = converted * gramsPerMl
	case dimension == units.DimensionVolume && unitDimension == units.DimensionMass:
		normalized = converted / gramsPerMl
	default:
		return nil
	}

	resolved := &resolvedUnit{
		NormalizedQuantity: normalized,
		NormalizedUnit:     units.StandardUnit[dimension],
	}
	return resolved
}

// Reads the unit as the name of one of the food's portions, and records which one
// it was so the item keeps the equivalent it was calculated with. Returns nil when
// the food has no portion with that name.
func resolveNamedPortion(food Food, quantity float64, unit string) *resolvedUnit {
	for _, portion := range food.Portions {
		if unit != portion.Name {
			continue
		}

		portionId := portion.Id
		portionEquivalent := portion.BaseEquivalent
		resolved := &resolvedUnit{
			NormalizedQuantity: quantity * portion.BaseEquivalent,
			NormalizedUnit:     units.StandardUnit[units.Dimension(food.MeasurementType)],
			PortionId:          &portionId,
			PortionEquivalent:  &portionEquivalent,
		}
		return resolved
	}
	return nil
}
