package ports

import (
	"github.com/google/uuid"
)

// What the food service needs to know about meals: whether a food is in use.
type FoodUsage interface {
	CountMealItems(foodId uuid.UUID) (int64, error)
}
