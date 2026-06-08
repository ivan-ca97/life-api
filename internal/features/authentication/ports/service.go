package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/authentication/domain"
)

type AuthenticationService interface {
	CreateSession(userId uuid.UUID) (*domain.Session, error)
	Validate(sessionId uuid.UUID) (*domain.Session, error)
	Logout(sessionId uuid.UUID) error
}
