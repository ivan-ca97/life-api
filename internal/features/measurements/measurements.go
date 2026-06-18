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
	repo := repository.NewMeasurementRepository(db)
	svc := service.NewMeasurementService(repo)
	authorizedSvc := service.NewAuthorizedMeasurementService(svc, authorizer)
	h := handler.NewMeasurementHandler(authorizedSvc)

	return &measurementsFeature{
		handler:      h,
		errorHandler: errorHandler,
	}
}
