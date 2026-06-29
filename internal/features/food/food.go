package food

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/food/handler"
	"github.com/ivan-ca97/life/internal/features/food/ports"
	"github.com/ivan-ca97/life/internal/features/food/repository"
	"github.com/ivan-ca97/life/internal/features/food/service"
)

type foodFeature struct {
	foodService  ports.FoodService
	foodHandler  handler.FoodHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewFoodFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *foodFeature {
	foodRepository := repository.NewFoodRepository(db)
	foodService := service.NewFoodService(foodRepository)
	authorizedService := service.NewAuthorizedFoodService(foodService, authorizer)
	foodHandler := handler.NewFoodHandler(authorizedService)

	return &foodFeature{
		foodService:  foodService,
		foodHandler:  foodHandler,
		errorHandler: errorHandler,
	}
}

// FoodService exposes the base food service so applications (e.g. the meal AI
// assistant) can search a user's catalog. Authorization is enforced by the
// caller for the owning user.
func (f *foodFeature) FoodService() ports.FoodService { return f.foodService }
