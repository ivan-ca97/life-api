package ports

import (
	"time"

	"github.com/google/uuid"
)

type AuthenticationResult struct {
	UserId    uuid.UUID
	Token     uuid.UUID
	ExpiresAt time.Time
}

type AuthenticationUseCase interface {
	Login(email, password string) (*AuthenticationResult, error)
	Register(email, password string) (*AuthenticationResult, error)
	LoginWithGoogle(idToken string) (*AuthenticationResult, error)
	Logout(sessionId uuid.UUID) error
}

type RoleAssigner interface {
	AssignRoleByName(userId uuid.UUID, roleName string) error
}
