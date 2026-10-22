package domain

import (
	"fmt"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

var (
	ErrFoodNotFound      = cerr.NewNotFoundError("food")
	ErrFoodAlreadyExists = cerr.NewConflictError("a food with this name already exists")

	ErrPortionNotFound      = cerr.NewNotFoundError("portion")
	ErrPortionAlreadyExists = cerr.NewConflictError("the food already has a portion with this name")
	ErrPortionNameRequired  = cerr.NewBadRequestError("portion name is required")
)

func NewInvalidPortionEquivalentError(name string) error {
	return cerr.NewBadRequestError(fmt.Sprintf("portion '%s' base_equivalent must be positive", name))
}
