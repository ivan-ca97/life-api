package service

import (
	"time"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/measurements/domain"
	"github.com/ivan-ca97/life/internal/features/measurements/ports"
)

type measurementService struct {
	repository ports.MeasurementRepository
}

var _ ports.MeasurementService = (*measurementService)(nil)

func NewMeasurementService(repository ports.MeasurementRepository) *measurementService {
	return &measurementService{repository: repository}
}

func (s *measurementService) Upsert(userId uuid.UUID, date time.Time, measureType string, params ports.UpsertParams) (*domain.BodyMeasurement, error) {
	if measureType == "" {
		return nil, cerr.NewBadRequestError("type is required")
	}
	if params.Value <= 0 {
		return nil, cerr.NewBadRequestError("value must be greater than 0")
	}
	m := &domain.BodyMeasurement{
		Id:     uuid.New(),
		UserId: userId,
		Date:   date,
		Type:   measureType,
		Value:  params.Value,
		Notes:  params.Notes,
	}
	if err := s.repository.Upsert(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *measurementService) GetByDate(userId uuid.UUID, date time.Time, measureType string) (*domain.BodyMeasurement, error) {
	return s.repository.FindByDate(userId, date, measureType)
}

func (s *measurementService) List(userId uuid.UUID, params ports.ListParams) ([]domain.BodyMeasurement, error) {
	return s.repository.List(userId, params)
}

func (s *measurementService) Delete(userId uuid.UUID, date time.Time, measureType string) error {
	return s.repository.Delete(userId, date, measureType)
}
