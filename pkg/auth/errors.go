package auth

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrNoActor   = cerr.NewUnauthorizedError("unauthenticated")
	ErrForbidden = cerr.NewForbiddenError("forbidden")
)
