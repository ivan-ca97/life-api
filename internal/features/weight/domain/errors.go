package domain

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrWeightEntryNotFound = cerr.NewNotFoundError("weight entry")
	ErrWeightEntryConflict = cerr.NewConflictError("weight entry already exists for this date")
)
