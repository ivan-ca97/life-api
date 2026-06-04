package domain

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrGoalNotFound = cerr.NewNotFoundError("goal")
)
