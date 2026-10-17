package domain

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrFoodNotFound      = cerr.NewNotFoundError("food")
	ErrFoodAlreadyExists = cerr.NewConflictError("a food with this name already exists")
)
