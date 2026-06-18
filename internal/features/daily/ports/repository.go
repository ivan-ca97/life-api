package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type SummaryRepository interface {
	GetDailySummary(userId uuid.UUID, date time.Time) (*domain.DailySummary, error)
	GetDailySummaryRange(userId uuid.UUID, from, to time.Time) ([]domain.DailySummary, error)
	GetDailyCheck(userId uuid.UUID, date time.Time) (*domain.DailyCheck, error)
}
