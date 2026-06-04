package meal

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *mealFeature) PublicRoutes(_ chi.Router) {}

func (f *mealFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/meals", endpoint.JSON(f.errorHandler, f.mealHandler.Create))
	r.Get("/meals", endpoint.JSON(f.errorHandler, f.mealHandler.List))
	r.Post("/meals/preview", endpoint.JSON(f.errorHandler, f.mealHandler.PreviewNutrition))
	r.Get("/meals/types", endpoint.JSON(f.errorHandler, f.mealHandler.ListTypes))
	r.Get("/meals/{id}", endpoint.JSON(f.errorHandler, f.mealHandler.GetById))
	r.Patch("/meals/{id}", endpoint.JSON(f.errorHandler, f.mealHandler.Update))
	r.Delete("/meals/{id}", endpoint.JSON(f.errorHandler, f.mealHandler.Delete))
}
