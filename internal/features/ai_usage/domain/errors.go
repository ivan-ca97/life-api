package domain

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrTierNotFound  = cerr.NewNotFoundError("ai tier")
	ErrQuotaExceeded = cerr.NewTooManyRequestsError("monthly AI usage limit reached")
	ErrTierNameTaken = cerr.NewConflictError("an ai tier with that name already exists")
)
