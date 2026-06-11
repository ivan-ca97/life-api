package exercise

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/features/exercise/handler"
	"github.com/ivan-ca97/life/internal/features/exercise/ports"
	"github.com/ivan-ca97/life/internal/features/exercise/repository"
	"github.com/ivan-ca97/life/internal/features/exercise/service"
)

type exerciseFeature struct {
	exerciseHandler    handler.ExerciseHandler
	exerciseService    ports.ExerciseService
	exerciseRepository ports.ExerciseRepository
	errorHandler       http_errors.HttpErrorHandler
}

func NewExerciseFeature(db *gorm.DB, authorizer auth.AuthorizationService, closureChecker dayclosure.DayClosureChecker, errorHandler http_errors.HttpErrorHandler) *exerciseFeature {
	exerciseRepository := repository.NewExerciseRepository(db)
	exerciseService := service.NewExerciseService(exerciseRepository, closureChecker)
	authorizedService := service.NewAuthorizedExerciseService(exerciseService, authorizer)
	exerciseHandler := handler.NewExerciseHandler(authorizedService)

	return &exerciseFeature{
		exerciseHandler:    exerciseHandler,
		exerciseService:    exerciseService,
		exerciseRepository: exerciseRepository,
		errorHandler:       errorHandler,
	}
}

func (f *exerciseFeature) ExerciseService() ports.ExerciseService {
	return f.exerciseService
}

func (f *exerciseFeature) Repository() ports.ExerciseRepository {
	return f.exerciseRepository
}
