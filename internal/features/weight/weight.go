package weight

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/features/weight/handler"
	"github.com/ivan-ca97/life/internal/features/weight/repository"
	"github.com/ivan-ca97/life/internal/features/weight/service"
)

type weightFeature struct {
	weightEntryHandler handler.WeightEntryHandler
	errorHandler       http_errors.HttpErrorHandler
}

func NewWeightFeature(db *gorm.DB, authorizer auth.AuthorizationService, closureChecker dayclosure.DayClosureChecker, errorHandler http_errors.HttpErrorHandler) *weightFeature {
	weightEntryRepository := repository.NewWeightEntryRepository(db)
	weightEntryService := service.NewWeightEntryService(weightEntryRepository, closureChecker)
	authorizedService := service.NewAuthorizedWeightEntryService(weightEntryService, authorizer)
	weightEntryHandler := handler.NewWeightEntryHandler(authorizedService)

	return &weightFeature{
		weightEntryHandler: weightEntryHandler,
		errorHandler:       errorHandler,
	}
}
