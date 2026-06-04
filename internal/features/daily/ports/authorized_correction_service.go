package ports

import (
	"context"
	"time"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type AuthorizedCorrectionService interface {
	GetCorrection(ctx context.Context, date time.Time) (*domain.Correction, error)
	UpsertCorrection(ctx context.Context, correction *domain.Correction) error
	DeleteCorrection(ctx context.Context, date time.Time) error
}
