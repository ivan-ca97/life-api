package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type AuthorizedSummaryService interface {
	GetSummary(ctx context.Context, ownerId uuid.UUID, date time.Time) (*domain.DailySummary, error)
	GetSummaryRange(ctx context.Context, ownerId uuid.UUID, from, to time.Time) ([]domain.DailySummary, error)
}
