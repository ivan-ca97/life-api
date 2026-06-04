package auth

import "context"

type AuthorizationService interface {
	Require(ctx context.Context, permission string) error
	RequireOn(ctx context.Context, permission string, resource any) error
}
