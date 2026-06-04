package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type SummaryService interface {
	GetSummary(userId uuid.UUID, date time.Time) (*domain.DailySummary, error)
}
