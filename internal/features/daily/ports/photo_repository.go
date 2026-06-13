package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type PhotoRepository interface {
	Create(photo *domain.DailyPhoto) error
	FindById(id, userId uuid.UUID) (*domain.DailyPhoto, error)
	ListByDate(userId uuid.UUID, date time.Time) ([]domain.DailyPhoto, error)
	Update(id, userId uuid.UUID, name *string, isPrimary *bool) (*domain.DailyPhoto, error)
	UnsetPrimary(userId uuid.UUID, date time.Time) error
	Delete(id, userId uuid.UUID) error
}
