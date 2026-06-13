package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type CreatePhotoParams struct {
	Date      time.Time
	Url       string
	Name      string
	IsPrimary bool
}

type UpdatePhotoParams struct {
	Name      *string
	IsPrimary *bool
}

type PhotoService interface {
	Create(userId uuid.UUID, params CreatePhotoParams) (*domain.DailyPhoto, error)
	List(userId uuid.UUID, date time.Time) ([]domain.DailyPhoto, error)
	Update(id, userId uuid.UUID, params UpdatePhotoParams) (*domain.DailyPhoto, error)
	Delete(id, userId uuid.UUID) error
}
