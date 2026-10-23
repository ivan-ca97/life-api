// Optional data of a new food, passed to domain.NewFood. It covers only the
// fields NewFood does not require.
package food_option

import (
	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
	"github.com/ivan-ca97/life/internal/features/nutrition/domain/food_change"
)

func Photo(url *string) domain.FoodOption {
	return domain.FoodOption(food_change.Photo(url))
}

func Calories(amount *float64) domain.FoodOption {
	return domain.FoodOption(food_change.Calories(amount))
}

func ProteinGrams(amount *float64) domain.FoodOption {
	return domain.FoodOption(food_change.ProteinGrams(amount))
}

func CarbsGrams(amount *float64) domain.FoodOption {
	return domain.FoodOption(food_change.CarbsGrams(amount))
}

func FatGrams(amount *float64) domain.FoodOption {
	return domain.FoodOption(food_change.FatGrams(amount))
}

func FiberGrams(amount *float64) domain.FoodOption {
	return domain.FoodOption(food_change.FiberGrams(amount))
}

func Density(density *domain.VolumeConversion) domain.FoodOption {
	return func(food *domain.Food) error {
		food.VolumeConversion = density
		return nil
	}
}

func Tags(tags *[]string) domain.FoodOption {
	return domain.FoodOption(food_change.Tags(tags))
}

func Ingredients(names *[]string) domain.FoodOption {
	return domain.FoodOption(food_change.Ingredients(names))
}
