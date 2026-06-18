package service

import (
	"time"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/steps/domain"
	"github.com/ivan-ca97/life/internal/features/steps/ports"
)

type stepsService struct {
	repository ports.StepsRepository
}

var _ ports.StepsService = (*stepsService)(nil)

func NewStepsService(repository ports.StepsRepository) *stepsService {
	return &stepsService{repository: repository}
}

func (s *stepsService) Upsert(userId uuid.UUID, date time.Time, params ports.UpsertParams) (*domain.DailySteps, error) {
	if params.Steps < 0 {
		return nil, cerr.NewBadRequestError("steps cannot be negative")
	}
	entry := &domain.DailySteps{
		UserId: userId,
		Date:   date,
		Steps:  params.Steps,
		Source: params.Source,
	}
	if err := s.repository.Upsert(entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *stepsService) GetByDate(userId uuid.UUID, date time.Time) (*domain.DailySteps, error) {
	return s.repository.FindByDate(userId, date)
}

func (s *stepsService) List(userId uuid.UUID, params ports.ListParams) ([]domain.DailySteps, error) {
	return s.repository.List(userId, params)
}

func (s *stepsService) Delete(userId uuid.UUID, date time.Time) error {
	return s.repository.Delete(userId, date)
}
