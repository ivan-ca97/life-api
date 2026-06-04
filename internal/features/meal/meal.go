package meal

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/meal/handler"
	"github.com/ivan-ca97/life/internal/features/meal/repository"
	"github.com/ivan-ca97/life/internal/features/meal/service"
)

type mealFeature struct {
	mealHandler  handler.MealHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewMealFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *mealFeature {
	mealRepository := repository.NewMealRepository(db)
	foodLookup := repository.NewFoodLookup(db)
	mealService := service.NewMealService(mealRepository, foodLookup)
	authorizedService := service.NewAuthorizedMealService(mealService, authorizer)
	mealHandler := handler.NewMealHandler(authorizedService)

	return &mealFeature{
		mealHandler:  mealHandler,
		errorHandler: errorHandler,
	}
}
