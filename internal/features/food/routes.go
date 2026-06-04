package food

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *foodFeature) PublicRoutes(_ chi.Router) {}

func (f *foodFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/foods", endpoint.JSON(f.errorHandler, f.foodHandler.Create))
	r.Get("/foods", endpoint.JSON(f.errorHandler, f.foodHandler.List))
	r.Get("/foods/units", endpoint.JSON(f.errorHandler, f.foodHandler.ListUnits))
	r.Get("/foods/frequency", endpoint.JSON(f.errorHandler, f.foodHandler.Frequency))
	r.Get("/foods/ingredients", endpoint.JSON(f.errorHandler, f.foodHandler.ListIngredients))
	r.Get("/foods/ingredients/frequency", endpoint.JSON(f.errorHandler, f.foodHandler.IngredientFrequency))
	r.Get("/foods/{id}/units", endpoint.JSON(f.errorHandler, f.foodHandler.ListFoodUnits))
	r.Get("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.GetById))
	r.Patch("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.Update))
	r.Delete("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.Delete))
}
