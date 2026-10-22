package domain

import (
	"fmt"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

var (
	ErrMealNotFound       = cerr.NewNotFoundError("meal")
	ErrItemUnitAndPortion = cerr.NewBadRequestError("a meal item is measured in a unit or in a portion, not both")
)

func NewUnknownUnitError(foodId uuid.UUID, unit string) error {
	return cerr.NewBadRequestError(fmt.Sprintf("unknown unit '%s' for food %s", unit, foodId))
}

func NewUnknownPortionError(foodId uuid.UUID, portionId uuid.UUID) error {
	return cerr.NewBadRequestError(fmt.Sprintf("portion %s does not belong to food %s", portionId, foodId))
}
