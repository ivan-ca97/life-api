package ports

import (
	"context"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
)

type AuthorizedGoalService interface {
	GetCurrent(ctx context.Context) (*domain.Goal, error)
	Upsert(ctx context.Context, params UpsertParams) (*domain.Goal, error)
}
