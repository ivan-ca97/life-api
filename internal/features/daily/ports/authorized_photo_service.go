package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type AuthorizedPhotoService interface {
	Create(ctx context.Context, ownerId uuid.UUID, params CreatePhotoParams) (*domain.DailyPhoto, error)
	List(ctx context.Context, ownerId uuid.UUID, date time.Time) ([]domain.DailyPhoto, error)
	Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params UpdatePhotoParams) (*domain.DailyPhoto, error)
	Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error
}
