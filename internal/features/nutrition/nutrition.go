package nutrition

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/features/nutrition/handler"
	"github.com/ivan-ca97/life/internal/features/nutrition/ports"
	"github.com/ivan-ca97/life/internal/features/nutrition/repository"
	"github.com/ivan-ca97/life/internal/features/nutrition/service"
)

// nutritionFeature owns the food catalog and the meal log. They depend on each
// other: meals are computed from foods, and food statistics are measured over
// meals.
type nutritionFeature struct {
	foodService  ports.FoodService
	foodHandler  handler.FoodHandler
	mealHandler  handler.MealHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewNutritionFeature(
	db *gorm.DB,
	authorizer auth.AuthorizationService,
	closureChecker dayclosure.DayClosureChecker,
	errorHandler http_errors.HttpErrorHandler,
) *nutritionFeature {
	foodRepository := repository.NewFoodRepository(db)
	foodService := service.NewFoodService(foodRepository)
	authorizedFoodService := service.NewAuthorizedFoodService(foodService, authorizer)
	foodHandler := handler.NewFoodHandler(authorizedFoodService)

	mealRepository := repository.NewMealRepository(db)
	foodLookup := repository.NewFoodLookup(db)
	mealService := service.NewMealService(mealRepository, foodLookup, closureChecker)
	authorizedMealService := service.NewAuthorizedMealService(mealService, authorizer)
	mealHandler := handler.NewMealHandler(authorizedMealService)

	return &nutritionFeature{
		foodService:  foodService,
		foodHandler:  foodHandler,
		mealHandler:  mealHandler,
		errorHandler: errorHandler,
	}
}

// FoodService exposes the base food service so applications (e.g. the meal AI
// assistant) can search a user's catalog. Authorization is enforced by the
// caller for the owning user.
func (f *nutritionFeature) FoodService() ports.FoodService {
	return f.foodService
}
