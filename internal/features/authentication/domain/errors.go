package domain

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrSessionNotFound    = cerr.NewUnauthorizedError("session not found")
	ErrSessionExpired     = cerr.NewUnauthorizedError("session expired")
	ErrInvalidGoogleToken = cerr.NewBadRequestError("invalid google token")
)
