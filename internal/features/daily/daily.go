package daily

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/features/daily/handler"
	"github.com/ivan-ca97/life/internal/features/daily/repository"
	"github.com/ivan-ca97/life/internal/features/daily/service"
)

type dailyFeature struct {
	summaryHandler    handler.SummaryHandler
	correctionHandler handler.CorrectionHandler
	closureHandler    handler.DayClosureHandler
	closureChecker    dayclosure.DayClosureChecker
	errorHandler      http_errors.HttpErrorHandler
}

func NewDailyFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *dailyFeature {
	summaryRepository := repository.NewSummaryRepository(db)
	summaryService := service.NewSummaryService(summaryRepository)
	authorizedSummaryService := service.NewAuthorizedSummaryService(summaryService, authorizer)
	summaryHandler := handler.NewSummaryHandler(authorizedSummaryService)

	correctionRepository := repository.NewCorrectionRepository(db)
	correctionService := service.NewCorrectionService(correctionRepository)
	authorizedCorrectionService := service.NewAuthorizedCorrectionService(correctionService, authorizer)
	correctionHandler := handler.NewCorrectionHandler(authorizedCorrectionService)

	dayClosureRepository := repository.NewDayClosureRepository(db)
	dayClosureService := service.NewDayClosureService(dayClosureRepository)
	authorizedDayClosureService := service.NewAuthorizedDayClosureService(dayClosureService, authorizer)
	closureHandler := handler.NewDayClosureHandler(authorizedDayClosureService)

	return &dailyFeature{
		summaryHandler:    summaryHandler,
		correctionHandler: correctionHandler,
		closureHandler:    closureHandler,
		closureChecker:    dayClosureService,
		errorHandler:      errorHandler,
	}
}

func (f *dailyFeature) DayClosureChecker() dayclosure.DayClosureChecker {
	return f.closureChecker
}
