package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/auth/domain"
)

type AuthService interface {
	Login(email, password string) (*domain.Session, error)
	CreateSession(userId uuid.UUID) (*domain.Session, error)
	Logout(sessionId uuid.UUID) error
	Validate(sessionId uuid.UUID) (*domain.Session, error)
}

type RoleAssigner interface {
	AssignRoleByName(userId uuid.UUID, roleName string) error
}
