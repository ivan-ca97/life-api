package nutrition

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/features/nutrition/food"
	foodPorts "github.com/ivan-ca97/life/internal/features/nutrition/food/ports"
	"github.com/ivan-ca97/life/internal/features/nutrition/meal"
)

// foodModule is what the feature needs from the food catalog: its routes and
// the service other modules use to read foods.
type foodModule interface {
	PublicRoutes(r chi.Router)
	GlobalRoutes(r chi.Router)
	ProtectedRoutes(r chi.Router)
	FoodService() foodPorts.FoodService
}

type mealModule interface {
	PublicRoutes(r chi.Router)
	ProtectedRoutes(r chi.Router)
}

// nutritionFeature groups the food catalog and the meal log, which depend on
// each other: meals are computed from foods and food statistics are measured
// over meals. To the rest of the system it is a single feature.
type nutritionFeature struct {
	foods foodModule
	meals mealModule
}

func NewNutritionFeature(
	db *gorm.DB,
	authorizer auth.AuthorizationService,
	closureChecker dayclosure.DayClosureChecker,
	errorHandler http_errors.HttpErrorHandler,
) *nutritionFeature {
	foods := food.NewFoodFeature(db, authorizer, errorHandler)
	meals := meal.NewMealFeature(db, authorizer, closureChecker, errorHandler)

	return &nutritionFeature{
		foods: foods,
		meals: meals,
	}
}

// FoodService exposes the base food service so applications (e.g. the meal AI
// assistant) can search a user's catalog. Authorization is enforced by the
// caller for the owning user.
func (f *nutritionFeature) FoodService() foodPorts.FoodService {
	return f.foods.FoodService()
}
