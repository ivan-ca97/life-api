package domain

import cerr "github.com/ivan-ca97/life/pkg/custom_error"

var (
	ErrNoInput       = cerr.NewBadRequestError("provide at least one photo or some instructions")
	ErrTooManyPhotos = cerr.NewBadRequestError("too many photos for a single estimation")
	ErrAIUnavailable = cerr.NewInternalError("ai estimation failed", nil)
)
