package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/authentication/domain"
)

type SessionRepository interface {
	Create(session *domain.Session) error
	FindById(id uuid.UUID) (*domain.Session, error)
	Delete(id uuid.UUID) error
	DeleteExpiredIfAbove(maxSessions int) error
}
