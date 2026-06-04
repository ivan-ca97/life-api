package ports

import (
	"context"
	"time"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type AuthorizedSummaryService interface {
	GetSummary(ctx context.Context, date time.Time) (*domain.DailySummary, error)
	GetSummaryRange(ctx context.Context, from, to time.Time) ([]domain.DailySummary, error)
}
