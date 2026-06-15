package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthorizationService interface {
	// Authorize checks if the current actor can perform the given permission on the specified owner's data.
	Authorize(ctx context.Context, ownerId uuid.UUID, permission string) error
	// AuthorizeAdmin returns nil only if the current actor has the admin role.
	AuthorizeAdmin(ctx context.Context) error
}
