package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/measurements/domain"
	"github.com/ivan-ca97/life/internal/features/measurements/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedMeasurementService struct {
	base       ports.MeasurementService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedMeasurementService = (*authorizedMeasurementService)(nil)

func NewAuthorizedMeasurementService(base ports.MeasurementService, authorizer auth.AuthorizationService) *authorizedMeasurementService {
	return &authorizedMeasurementService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedMeasurementService) Upsert(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string, params ports.UpsertParams) (*domain.BodyMeasurement, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MeasurementsWrite)
	if err != nil {
		return nil, err
	}
	measurement, err := s.base.Upsert(ownerId, date, measureType, params)
	if err != nil {
		return nil, err
	}
	return measurement, nil
}

func (s *authorizedMeasurementService) GetByDate(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string) (*domain.BodyMeasurement, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MeasurementsRead)
	if err != nil {
		return nil, err
	}
	measurement, err := s.base.GetByDate(ownerId, date, measureType)
	if err != nil {
		return nil, err
	}
	return measurement, nil
}

func (s *authorizedMeasurementService) List(ctx context.Context, ownerId uuid.UUID, params ports.ListParams) ([]domain.BodyMeasurement, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MeasurementsRead)
	if err != nil {
		return nil, err
	}
	measurements, err := s.base.List(ownerId, params)
	if err != nil {
		return nil, err
	}
	return measurements, nil
}

func (s *authorizedMeasurementService) Delete(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.MeasurementsWrite)
	if err != nil {
		return err
	}
	err = s.base.Delete(ownerId, date, measureType)
	if err != nil {
		return err
	}
	return nil
}
