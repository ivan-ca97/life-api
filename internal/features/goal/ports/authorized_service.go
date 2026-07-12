package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
)

type AuthorizedGoalService interface {
	GetCurrent(ctx context.Context, ownerId uuid.UUID) (*domain.Goal, error)
	Upsert(ctx context.Context, ownerId uuid.UUID, params UpsertParams) (*domain.Goal, error)
	GetProgress(ctx context.Context, ownerId uuid.UUID, from, to time.Time) (*domain.GoalProgress, error)
}
