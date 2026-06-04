package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type correctionService struct {
	repository ports.CorrectionRepository
}

var _ ports.CorrectionService = (*correctionService)(nil)

func NewCorrectionService(repository ports.CorrectionRepository) *correctionService {
	return &correctionService{repository: repository}
}

func (s *correctionService) GetCorrection(userId uuid.UUID, date time.Time) (*domain.Correction, error) {
	return s.repository.GetCorrection(userId, date)
}

func (s *correctionService) UpsertCorrection(userId uuid.UUID, correction *domain.Correction) error {
	return s.repository.UpsertCorrection(userId, correction)
}

func (s *correctionService) DeleteCorrection(userId uuid.UUID, date time.Time) error {
	return s.repository.DeleteCorrection(userId, date)
}
