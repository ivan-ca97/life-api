package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type AuthorizedCorrectionService interface {
	GetCorrection(ctx context.Context, ownerId uuid.UUID, date time.Time) (*domain.Correction, error)
	UpsertCorrection(ctx context.Context, ownerId uuid.UUID, correction *domain.Correction) error
	DeleteCorrection(ctx context.Context, ownerId uuid.UUID, date time.Time) error
}
