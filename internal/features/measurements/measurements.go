package measurements

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/measurements/handler"
	"github.com/ivan-ca97/life/internal/features/measurements/repository"
	"github.com/ivan-ca97/life/internal/features/measurements/service"
)

type measurementsFeature struct {
	handler      handler.MeasurementHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewMeasurementsFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *measurementsFeature {
	measurementRepository := repository.NewMeasurementRepository(db)
	measurementService := service.NewMeasurementService(measurementRepository)
	authorizedService := service.NewAuthorizedMeasurementService(measurementService, authorizer)
	measurementHandler := handler.NewMeasurementHandler(authorizedService)

	return &measurementsFeature{
		handler:      measurementHandler,
		errorHandler: errorHandler,
	}
}
