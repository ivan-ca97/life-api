package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/steps/domain"
)

type ListParams struct {
	From *time.Time
	To   *time.Time
}

type StepsRepository interface {
	Upsert(entry *domain.DailySteps) error
	FindByDate(userId uuid.UUID, date time.Time) (*domain.DailySteps, error)
	List(userId uuid.UUID, params ListParams) ([]domain.DailySteps, error)
	Delete(userId uuid.UUID, date time.Time) error
}
