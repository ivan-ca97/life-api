package meal_ai

import (
	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"

	foodPorts "github.com/ivan-ca97/life/internal/features/food/ports"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/handler"
	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
	"github.com/ivan-ca97/life/internal/applications/meal_ai/use_case"
)

type MealAIApplication struct {
	handler      handler.MealAIHandler
	errorHandler http_errors.HttpErrorHandler
}

// NewMealAIApplication wires the meal estimation feature. client is any LLM
// provider (e.g. the openai adapter); quota, logger and pricer come from the
// ai_usage feature.
func NewMealAIApplication(
	foodService foodPorts.FoodService,
	client llm.Client,
	quota ports.QuotaGuard,
	logger ports.InteractionLogger,
	pricer ports.Pricer,
	authorizer auth.AuthorizationService,
	errorHandler http_errors.HttpErrorHandler,
) *MealAIApplication {
	foodSearch := &foodSearchAdapter{foodService: foodService}
	imageFetcher := newHTTPImageFetcher()
	estimationUseCase := use_case.NewMealEstimationUseCase(client, foodSearch, imageFetcher, quota, logger, pricer, authorizer)
	mealAIHandler := handler.NewMealAIHandler(estimationUseCase)

	return &MealAIApplication{
		handler:      mealAIHandler,
		errorHandler: errorHandler,
	}
}
