package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
)

type CorrectionService interface {
	GetCorrection(userId uuid.UUID, date time.Time) (*domain.Correction, error)
	UpsertCorrection(userId uuid.UUID, correction *domain.Correction) error
	DeleteCorrection(userId uuid.UUID, date time.Time) error
}
