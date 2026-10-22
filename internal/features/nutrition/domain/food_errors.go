package domain

import (
	"fmt"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

var (
	ErrFoodNotFound      = cerr.NewNotFoundError("food")
	ErrFoodAlreadyExists = cerr.NewConflictError("a food with this name already exists")
	ErrFoodNameRequired  = cerr.NewBadRequestError("name is required")

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
