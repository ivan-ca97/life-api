// Changes to a food, passed to Food.Update. Each one sets a
// single field from an optional value: a nil value changes nothing, so the
// optional fields of a partial update map to changes one to one.
package food_change

import (
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
)

func Name(name *string) domain.FoodChange {
	return func(food *domain.Food) error {
		if name != nil {
			food.Name = *name
		}
		return nil
	}
}

func Photo(url *string) domain.FoodChange {
	return func(food *domain.Food) error {
		if url != nil {
			food.PhotoUrl = *url
		}
		return nil
	}
}

func Calories(amount *float64) domain.FoodChange {
	return func(food *domain.Food) error {
		if amount != nil {
			food.DefaultCalories = amount
		}
		return nil
	}
}

func ProteinGrams(amount *float64) domain.FoodChange {
	return func(food *domain.Food) error {
		if amount != nil {
			food.DefaultProteinGrams = amount
		}
		return nil
	}
}

func CarbsGrams(amount *float64) domain.FoodChange {
	return func(food *domain.Food) error {
		if amount != nil {
			food.DefaultCarbsGrams = amount
		}
		return nil
	}
}

func FatGrams(amount *float64) domain.FoodChange {
	return func(food *domain.Food) error {
		if amount != nil {
			food.DefaultFatGrams = amount
		}
		return nil
	}
}

func FiberGrams(amount *float64) domain.FoodChange {
	return func(food *domain.Food) error {
		if amount != nil {
			food.DefaultFiberGrams = amount
		}
		return nil
	}
}

func BaseQuantity(quantity *float64) domain.FoodChange {
	return func(food *domain.Food) error {
		if quantity != nil {
			food.BaseQuantity = *quantity
		}
		return nil
	}
}

func BaseUnit(name *string) domain.FoodChange {
	return func(food *domain.Food) error {
		if name == nil {
			return nil
		}
		unit, err := domain.ParseBaseUnit(*name)
		if err != nil {
			return err
		}
		food.BaseUnit = *unit
		return nil
	}
}

// A density that comes without a value removes it.
func Density(density types.Optional[domain.VolumeConversion]) domain.FoodChange {
	return func(food *domain.Food) error {
		if density.Set {
			food.VolumeConversion = density.Value
		}
		return nil
	}
}

func Tags(tags *[]string) domain.FoodChange {
	return func(food *domain.Food) error {
		if tags != nil {
			food.Tags = *tags
		}
		return nil
	}
}

func Ingredients(names *[]string) domain.FoodChange {
	return func(food *domain.Food) error {
		if names == nil {
			return nil
		}
		ingredients := make([]domain.Ingredient, len(*names))
		for i, name := range *names {
			ingredients[i] = domain.Ingredient{
				Name: name,
			}
		}
		food.Ingredients = ingredients
		return nil
	}
}
