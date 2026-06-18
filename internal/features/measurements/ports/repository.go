package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/measurements/domain"
)

type ListParams struct {
	From *time.Time
	To   *time.Time
	Type *string
}

type MeasurementRepository interface {
	Upsert(m *domain.BodyMeasurement) error
	FindByDate(userId uuid.UUID, date time.Time, measureType string) (*domain.BodyMeasurement, error)
	List(userId uuid.UUID, params ListParams) ([]domain.BodyMeasurement, error)
	Delete(userId uuid.UUID, date time.Time, measureType string) error
}
