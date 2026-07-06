package ai_usage

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/ai_usage/handler"
	"github.com/ivan-ca97/life/internal/features/ai_usage/ports"
	"github.com/ivan-ca97/life/internal/features/ai_usage/repository"
	"github.com/ivan-ca97/life/internal/features/ai_usage/service"
)

type aiUsageFeature struct {
	service      ports.Service
	handler      handler.AiUsageHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewAiUsageFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *aiUsageFeature {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	authorizedService := service.NewAuthorizedService(svc, authorizer)
	aiUsageHandler := handler.NewAiUsageHandler(authorizedService)

	return &aiUsageFeature{
		service:      svc,
		handler:      aiUsageHandler,
		errorHandler: errorHandler,
	}
}

// QuotaGuard exposes the spend-limit checks the meal AI feature consumes.
func (f *aiUsageFeature) QuotaGuard() ports.QuotaGuard { return f.service }

// InteractionLogger exposes the interaction log the meal AI feature writes to.
func (f *aiUsageFeature) InteractionLogger() ports.InteractionLogger { return f.service }
