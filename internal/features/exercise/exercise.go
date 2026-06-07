package exercise

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/features/exercise/handler"
	"github.com/ivan-ca97/life/internal/features/exercise/repository"
	"github.com/ivan-ca97/life/internal/features/exercise/service"
)

type exerciseFeature struct {
	exerciseHandler handler.ExerciseHandler
	errorHandler    http_errors.HttpErrorHandler
}

func NewExerciseFeature(db *gorm.DB, authorizer auth.AuthorizationService, closureChecker dayclosure.DayClosureChecker, errorHandler http_errors.HttpErrorHandler) *exerciseFeature {
	exerciseRepository := repository.NewExerciseRepository(db)
	exerciseService := service.NewExerciseService(exerciseRepository)
	authorizedService := service.NewAuthorizedExerciseService(exerciseService, authorizer, closureChecker)
	exerciseHandler := handler.NewExerciseHandler(authorizedService)

	return &exerciseFeature{
		exerciseHandler: exerciseHandler,
		errorHandler:    errorHandler,
	}
}
