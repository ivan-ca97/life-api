package health_connect_import

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	exercisePorts "github.com/ivan-ca97/life/internal/features/exercise/ports"
	weightPorts "github.com/ivan-ca97/life/internal/features/weight/ports"

	"github.com/ivan-ca97/life/internal/applications/health_connect_import/handler"
	"github.com/ivan-ca97/life/internal/applications/health_connect_import/repository"
	"github.com/ivan-ca97/life/internal/applications/health_connect_import/use_case"
)

type healthConnectImportApplication struct {
	importHandler handler.HealthConnectImportHandler
	errorHandler  http_errors.HttpErrorHandler
}

func NewHealthConnectImportApplication(
	db *gorm.DB,
	weightService weightPorts.WeightEntryService,
	weightRepository weightPorts.WeightEntryRepository,
	exerciseService exercisePorts.ExerciseService,
	exerciseRepository exercisePorts.ExerciseRepository,
	authorizer auth.AuthorizationService,
	errorHandler http_errors.HttpErrorHandler,
) *healthConnectImportApplication {
	rawStore := repository.NewExternalHealthRecordRepository(db)
	syncLogs := repository.NewSyncLogRepository(db)
	importUseCase := use_case.NewHealthConnectImportUseCase(
		weightService,
		weightRepository,
		exerciseService,
		exerciseRepository,
		rawStore,
		syncLogs,
		authorizer,
	)

	importHandler := handler.NewHealthConnectImportHandler(importUseCase, nil)

	return &healthConnectImportApplication{
		importHandler: importHandler,
		errorHandler:  errorHandler,
	}
}
