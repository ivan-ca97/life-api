package domain

import (
	"fmt"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

var (
	ErrFoodNotFound        = cerr.NewNotFoundError("food")
	ErrFoodAlreadyExists   = cerr.NewConflictError("a food with this name already exists")
	ErrFoodNameRequired    = cerr.NewBadRequestError("name is required")
	ErrInvalidBaseQuantity = cerr.NewBadRequestError("base_quantity must be positive")
	ErrNegativeNutrient    = cerr.NewBadRequestError("nutritional values cannot be negative")
	ErrInvalidDensity      = cerr.NewBadRequestError("grams_per_ml must be positive")

	ErrPortionNotFound      = cerr.NewNotFoundError("portion")
	ErrPortionAlreadyExists = cerr.NewConflictError("the food already has a portion with this name")
	ErrPortionNameRequired  = cerr.NewBadRequestError("portion name is required")
)

func NewInvalidBaseUnitError(unit string) error {
	return cerr.NewBadRequestError(fmt.Sprintf("'%s' is not a valid metric unit", unit))
}

func NewInvalidPortionEquivalentError(name string) error {
	return cerr.NewBadRequestError(fmt.Sprintf("portion '%s' base_equivalent must be positive", name))
}

func NewUnknownUnitError(foodId uuid.UUID, unit string) error {
	return cerr.NewBadRequestError(fmt.Sprintf("unknown unit '%s' for food %s", unit, foodId))
}

func NewUnknownPortionError(foodId uuid.UUID, portionId uuid.UUID) error {
	return cerr.NewBadRequestError(fmt.Sprintf("portion %s does not belong to food %s", portionId, foodId))
}
