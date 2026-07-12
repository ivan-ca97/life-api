package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
)

type GoalRepository interface {
	FindByUserId(userId uuid.UUID) (*domain.Goal, error)
	Upsert(goal *domain.Goal) (*domain.Goal, error)
	GetProgress(userId uuid.UUID, from, to time.Time) (*domain.GoalProgress, error)
}
