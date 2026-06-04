package goal

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/goal/handler"
	"github.com/ivan-ca97/life/internal/features/goal/repository"
	"github.com/ivan-ca97/life/internal/features/goal/service"
)

type goalFeature struct {
	goalHandler  handler.GoalHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewGoalFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *goalFeature {
	goalRepository := repository.NewGoalRepository(db)
	goalService := service.NewGoalService(goalRepository)
	authorizedService := service.NewAuthorizedGoalService(goalService, authorizer)
	goalHandler := handler.NewGoalHandler(authorizedService)

	return &goalFeature{
		goalHandler:  goalHandler,
		errorHandler: errorHandler,
	}
}
