package domain

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrExerciseNotFound    = cerr.NewNotFoundError("exercise")
	ErrInvalidExerciseType = cerr.NewBadRequestError("invalid exercise type")
)
