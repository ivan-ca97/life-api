package food

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/food/handler"
	"github.com/ivan-ca97/life/internal/features/food/repository"
	"github.com/ivan-ca97/life/internal/features/food/service"
)

type foodFeature struct {
	foodHandler  handler.FoodHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewFoodFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *foodFeature {
	foodRepository := repository.NewFoodRepository(db)
	foodService := service.NewFoodService(foodRepository)
	authorizedService := service.NewAuthorizedFoodService(foodService, authorizer)
	foodHandler := handler.NewFoodHandler(authorizedService)

	return &foodFeature{
		foodHandler:  foodHandler,
		errorHandler: errorHandler,
	}
}
