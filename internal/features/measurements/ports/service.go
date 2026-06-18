package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/measurements/domain"
)

type UpsertParams struct {
	Value float64
	Notes string
}

type MeasurementService interface {
	Upsert(userId uuid.UUID, date time.Time, measureType string, params UpsertParams) (*domain.BodyMeasurement, error)
	GetByDate(userId uuid.UUID, date time.Time, measureType string) (*domain.BodyMeasurement, error)
	List(userId uuid.UUID, params ListParams) ([]domain.BodyMeasurement, error)
	Delete(userId uuid.UUID, date time.Time, measureType string) error
}

type AuthorizedMeasurementService interface {
	Upsert(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string, params UpsertParams) (*domain.BodyMeasurement, error)
	GetByDate(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string) (*domain.BodyMeasurement, error)
	List(ctx context.Context, ownerId uuid.UUID, params ListParams) ([]domain.BodyMeasurement, error)
	Delete(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string) error
}
