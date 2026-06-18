package steps

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/steps/handler"
	"github.com/ivan-ca97/life/internal/features/steps/repository"
	"github.com/ivan-ca97/life/internal/features/steps/service"
)

type stepsFeature struct {
	handler      handler.StepsHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewStepsFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *stepsFeature {
	stepsRepository := repository.NewStepsRepository(db)
	weightLookup := repository.NewWeightLookup(db)
	stepsService := service.NewStepsService(stepsRepository)
	authorizedService := service.NewAuthorizedStepsService(stepsService, authorizer)
	stepsHandler := handler.NewStepsHandler(authorizedService, weightLookup)

	return &stepsFeature{
		handler:      stepsHandler,
		errorHandler: errorHandler,
	}
}
