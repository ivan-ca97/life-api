package hevy_import

import (
	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	exercisePorts "github.com/ivan-ca97/life/internal/features/exercise/ports"

	"github.com/ivan-ca97/life/internal/applications/hevy_import/handler"
	"github.com/ivan-ca97/life/internal/applications/hevy_import/use_case"
)

type hevyImportApplication struct {
	importHandler handler.HevyImportHandler
	errorHandler  http_errors.HttpErrorHandler
}

func NewHevyImportApplication(
	exerciseService exercisePorts.ExerciseService,
	exerciseRepository exercisePorts.ExerciseRepository,
	authorizer auth.AuthorizationService,
	errorHandler http_errors.HttpErrorHandler,
) *hevyImportApplication {
	importUseCase := use_case.NewHevyImportUseCase(exerciseService, exerciseRepository, authorizer)
	importHandler := handler.NewHevyImportHandler(importUseCase)

	return &hevyImportApplication{
		importHandler: importHandler,
		errorHandler:  errorHandler,
	}
}
