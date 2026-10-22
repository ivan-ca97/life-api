package domain

import (
	"github.com/google/uuid"
)

// Builds the item that results from eating quantity of food.
func NewMealItem(food Food, quantity float64, unit string) (*MealItem, error) {
	normalized, err := food.InStandardUnit(quantity, unit)
	if err != nil {
		return nil, err
	}
	item := newMealItem(food, quantity, unit, normalized)
	return item, nil
}

// Builds the item that results from eating count of one of the food's portions.
// The item keeps the portion's name and equivalent, so editing the portion later
// does not change what was eaten.
func NewMealItemOfPortion(food Food, count float64, portionId uuid.UUID) (*MealItem, error) {
	portion, err := food.Portion(portionId)
	if err != nil {
		return nil, err
	}
	item := newMealItem(food, count, portion.Name, food.PortionInStandardUnit(*portion, count))
	item.InputPortionId = &portion.Id
	item.InputPortionEquivalent = &portion.BaseEquivalent
	return item, nil
}

func newMealItem(food Food, inputQuantity float64, inputUnit string, normalizedQuantity float64) *MealItem {
	nutrition := food.NutritionFor(normalizedQuantity)
	item := &MealItem{
		FoodId:             food.Id,
		InputQuantity:      inputQuantity,
		InputUnit:          inputUnit,
		NormalizedQuantity: normalizedQuantity,
		NormalizedUnit:     food.StandardUnit(),
		Calories:           nutrition.Calories,
		ProteinGrams:       nutrition.ProteinGrams,
		CarbsGrams:         nutrition.CarbsGrams,
		FatGrams:           nutrition.FatGrams,
		FiberGrams:         nutrition.FiberGrams,
	}
	return item
}
