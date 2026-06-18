package domain

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrUserNotFound  = cerr.NewNotFoundError("user")
	ErrEmailTaken    = cerr.NewConflictError("email already in use")
	ErrUsernameTaken = cerr.NewConflictError("username already in use")
	ErrUserInactive  = cerr.NewUnauthorizedError("user account is inactive")
)
