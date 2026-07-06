package meal_ai

import (
	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	foodPorts "github.com/ivan-ca97/life/internal/features/food/ports"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/handler"
	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
	"github.com/ivan-ca97/life/internal/applications/meal_ai/use_case"
)

type MealAIApplication struct {
	handler      handler.MealAIHandler
	errorHandler http_errors.HttpErrorHandler
}

// NewMealAIApplication wires the meal estimation feature. completer is satisfied
// by *openai.Client and model is its model id (e.g. "gpt-4o"); quota and logger
// come from the ai_usage feature.
func NewMealAIApplication(
	foodService foodPorts.FoodService,
	quota ports.QuotaGuard,
	logger ports.InteractionLogger,
	completer ports.Completer,
	model string,
	authorizer auth.AuthorizationService,
	errorHandler http_errors.HttpErrorHandler,
) *MealAIApplication {
	foodSearch := &foodSearchAdapter{foodService: foodService}
	imageFetcher := newHTTPImageFetcher()
	estimationUseCase := use_case.NewMealEstimationUseCase(completer, foodSearch, imageFetcher, quota, logger, authorizer, model)
	mealAIHandler := handler.NewMealAIHandler(estimationUseCase)

	return &MealAIApplication{
		handler:      mealAIHandler,
		errorHandler: errorHandler,
	}
}
