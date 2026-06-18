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
	return &authorizedMeasurementService{base: base, authorizer: authorizer}
}

func (s *authorizedMeasurementService) Upsert(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string, params ports.UpsertParams) (*domain.BodyMeasurement, error) {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.MeasurementsWrite); err != nil {
		return nil, err
	}
	return s.base.Upsert(ownerId, date, measureType, params)
}

func (s *authorizedMeasurementService) GetByDate(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string) (*domain.BodyMeasurement, error) {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.MeasurementsRead); err != nil {
		return nil, err
	}
	return s.base.GetByDate(ownerId, date, measureType)
}

func (s *authorizedMeasurementService) List(ctx context.Context, ownerId uuid.UUID, params ports.ListParams) ([]domain.BodyMeasurement, error) {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.MeasurementsRead); err != nil {
		return nil, err
	}
	return s.base.List(ownerId, params)
}

func (s *authorizedMeasurementService) Delete(ctx context.Context, ownerId uuid.UUID, date time.Time, measureType string) error {
	if err := s.authorizer.Authorize(ctx, ownerId, permissions.MeasurementsWrite); err != nil {
		return err
	}
	return s.base.Delete(ownerId, date, measureType)
}
