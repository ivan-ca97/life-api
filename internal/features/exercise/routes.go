package exercise

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *exerciseFeature) PublicRoutes(_ chi.Router) {}

func (f *exerciseFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/exercises", endpoint.JSON(f.errorHandler, f.exerciseHandler.Create))
	r.Get("/exercises", endpoint.JSON(f.errorHandler, f.exerciseHandler.List))
	r.Get("/exercises/{id}", endpoint.JSON(f.errorHandler, f.exerciseHandler.GetById))
	r.Patch("/exercises/{id}", endpoint.JSON(f.errorHandler, f.exerciseHandler.Update))
	r.Delete("/exercises/{id}", endpoint.JSON(f.errorHandler, f.exerciseHandler.Delete))
}
