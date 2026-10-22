package nutrition

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *nutritionFeature) PublicRoutes(_ chi.Router) {}

// GlobalRoutes are the routes that are not scoped to a user.
func (f *nutritionFeature) GlobalRoutes(r chi.Router) {
	r.Get("/foods/community", endpoint.JSON(f.errorHandler, f.foodHandler.ListCommunity))
	r.Get("/foods/units", endpoint.JSON(f.errorHandler, f.foodHandler.ListUnits))
}

func (f *nutritionFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/foods", endpoint.JSON(f.errorHandler, f.foodHandler.Create))
	r.Get("/foods", endpoint.JSON(f.errorHandler, f.foodHandler.List))
	r.Get("/foods/frequency", endpoint.JSON(f.errorHandler, f.mealHandler.FoodFrequency))
	r.Get("/foods/ingredients", endpoint.JSON(f.errorHandler, f.foodHandler.ListIngredients))
	r.Get("/foods/ingredients/frequency", endpoint.JSON(f.errorHandler, f.mealHandler.IngredientFrequency))
	r.Post("/foods/{id}/copy", endpoint.JSON(f.errorHandler, f.foodHandler.Copy))
	r.Get("/foods/{id}/units", endpoint.JSON(f.errorHandler, f.foodHandler.ListFoodUnits))
	r.Get("/foods/{id}/impact", endpoint.JSON(f.errorHandler, f.mealHandler.FoodImpact))
	r.Get("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.GetById))
	r.Patch("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.Update))
	r.Delete("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.Delete))
	r.Post("/foods/{id}/portions", endpoint.JSON(f.errorHandler, f.foodHandler.AddPortion))
	r.Patch("/foods/{id}/portions/{portionId}", endpoint.JSON(f.errorHandler, f.foodHandler.UpdatePortion))
	r.Delete("/foods/{id}/portions/{portionId}", endpoint.JSON(f.errorHandler, f.foodHandler.DeletePortion))

	r.Post("/meals", endpoint.JSON(f.errorHandler, f.mealHandler.Create))
	r.Get("/meals", endpoint.JSON(f.errorHandler, f.mealHandler.List))
	r.Post("/meals/preview", endpoint.JSON(f.errorHandler, f.mealHandler.PreviewNutrition))
	r.Get("/meals/types", endpoint.JSON(f.errorHandler, f.mealHandler.ListTypes))
	r.Get("/meals/{id}", endpoint.JSON(f.errorHandler, f.mealHandler.GetById))
	r.Patch("/meals/{id}", endpoint.JSON(f.errorHandler, f.mealHandler.Update))
	r.Delete("/meals/{id}", endpoint.JSON(f.errorHandler, f.mealHandler.Delete))
}
