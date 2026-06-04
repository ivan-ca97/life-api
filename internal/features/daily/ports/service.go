package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type SummaryService interface {
	GetSummary(userId uuid.UUID, date time.Time) (*domain.DailySummary, error)
	GetSummaryRange(userId uuid.UUID, from, to time.Time) ([]domain.DailySummary, error)
}
