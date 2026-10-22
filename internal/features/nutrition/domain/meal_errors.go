package domain

import (
	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

var (
	ErrMealNotFound       = cerr.NewNotFoundError("meal")
	ErrItemUnitAndPortion = cerr.NewBadRequestError("a meal item is measured in a unit or in a portion, not both")
)
