package domain

type NutritionTotals struct {
	Calories     *float64
	ProteinGrams *float64
	CarbsGrams   *float64
	FatGrams     *float64
	FiberGrams   *float64
}

// Adds up the nutrition of the items.
func TotalNutrition(items []MealItem) NutritionTotals {
	var totals NutritionTotals
	for _, item := range items {
		totals.Calories = addNutrient(totals.Calories, item.Calories)
		totals.ProteinGrams = addNutrient(totals.ProteinGrams, item.ProteinGrams)
		totals.CarbsGrams = addNutrient(totals.CarbsGrams, item.CarbsGrams)
		totals.FatGrams = addNutrient(totals.FatGrams, item.FatGrams)
		totals.FiberGrams = addNutrient(totals.FiberGrams, item.FiberGrams)
	}
	return totals
}

func addNutrient(total *float64, amount *float64) *float64 {
	if amount == nil {
		return total
	}
	if total == nil {
		sum := *amount
		return &sum
	}
	sum := *total + *amount
	return &sum
}
