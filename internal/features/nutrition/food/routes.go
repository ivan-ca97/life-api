package food

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *foodFeature) PublicRoutes(_ chi.Router) {}

func (f *foodFeature) GlobalRoutes(r chi.Router) {
	r.Get("/foods/community", endpoint.JSON(f.errorHandler, f.foodHandler.ListCommunity))
	r.Get("/foods/units", endpoint.JSON(f.errorHandler, f.foodHandler.ListUnits))
}

func (f *foodFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/foods", endpoint.JSON(f.errorHandler, f.foodHandler.Create))
	r.Get("/foods", endpoint.JSON(f.errorHandler, f.foodHandler.List))
	r.Get("/foods/frequency", endpoint.JSON(f.errorHandler, f.foodHandler.Frequency))
	r.Get("/foods/ingredients", endpoint.JSON(f.errorHandler, f.foodHandler.ListIngredients))
	r.Get("/foods/ingredients/frequency", endpoint.JSON(f.errorHandler, f.foodHandler.IngredientFrequency))
	r.Post("/foods/{id}/copy", endpoint.JSON(f.errorHandler, f.foodHandler.Copy))
	r.Get("/foods/{id}/units", endpoint.JSON(f.errorHandler, f.foodHandler.ListFoodUnits))
	r.Get("/foods/{id}/impact", endpoint.JSON(f.errorHandler, f.foodHandler.Impact))
	r.Get("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.GetById))
	r.Patch("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.Update))
	r.Delete("/foods/{id}", endpoint.JSON(f.errorHandler, f.foodHandler.Delete))
}
