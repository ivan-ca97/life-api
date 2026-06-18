package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/steps/domain"
)

type UpsertParams struct {
	Steps  int
	Source string
}

type StepsService interface {
	Upsert(userId uuid.UUID, date time.Time, params UpsertParams) (*domain.DailySteps, error)
	GetByDate(userId uuid.UUID, date time.Time) (*domain.DailySteps, error)
	List(userId uuid.UUID, params ListParams) ([]domain.DailySteps, error)
	Delete(userId uuid.UUID, date time.Time) error
}

type AuthorizedStepsService interface {
	Upsert(ctx context.Context, userId uuid.UUID, date time.Time, params UpsertParams) (*domain.DailySteps, error)
	GetByDate(ctx context.Context, userId uuid.UUID, date time.Time) (*domain.DailySteps, error)
	List(ctx context.Context, userId uuid.UUID, params ListParams) ([]domain.DailySteps, error)
	Delete(ctx context.Context, userId uuid.UUID, date time.Time) error
}
