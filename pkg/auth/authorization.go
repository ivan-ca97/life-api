package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthorizationService interface {
	// Authorize checks if the current actor can perform the given permission on the specified owner's data.
	Authorize(ctx context.Context, ownerId uuid.UUID, permission string) error
}
