package domain

import (
	"fmt"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

var (
	ErrMealNotFound = cerr.NewNotFoundError("meal")
)

func NewUnknownUnitError(foodId uuid.UUID, unit string) error {
	return cerr.NewBadRequestError(fmt.Sprintf("unknown unit '%s' for food %s", unit, foodId))
}
