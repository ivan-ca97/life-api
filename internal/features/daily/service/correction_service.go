package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type correctionService struct {
	repository     ports.CorrectionRepository
	closureChecker dayclosure.DayClosureChecker
}

var _ ports.CorrectionService = (*correctionService)(nil)

func NewCorrectionService(repository ports.CorrectionRepository, closureChecker dayclosure.DayClosureChecker) *correctionService {
	return &correctionService{
		repository:     repository,
		closureChecker: closureChecker,
	}
}

func (s *correctionService) GetCorrection(userId uuid.UUID, date time.Time) (*domain.Correction, error) {
	return s.repository.GetCorrection(userId, date)
}

func (s *correctionService) UpsertCorrection(userId uuid.UUID, correction *domain.Correction) error {
	closed, err := s.closureChecker.IsClosed(userId, correction.Date)
	if err != nil {
		return err
	}
	if closed {
		return dayclosure.ErrDayClosed
	}

	return s.repository.UpsertCorrection(userId, correction)
}

func (s *correctionService) DeleteCorrection(userId uuid.UUID, date time.Time) error {
	closed, err := s.closureChecker.IsClosed(userId, date)
	if err != nil {
		return err
	}
	if closed {
		return dayclosure.ErrDayClosed
	}

	return s.repository.DeleteCorrection(userId, date)
}
