package weight

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *weightFeature) PublicRoutes(_ chi.Router) {}

func (f *weightFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/weight", endpoint.JSON(f.errorHandler, f.weightEntryHandler.Create))
	r.Get("/weight", endpoint.JSON(f.errorHandler, f.weightEntryHandler.List))
	r.Get("/weight/{id}", endpoint.JSON(f.errorHandler, f.weightEntryHandler.GetById))
	r.Patch("/weight/{id}", endpoint.JSON(f.errorHandler, f.weightEntryHandler.Update))
	r.Delete("/weight/{id}", endpoint.JSON(f.errorHandler, f.weightEntryHandler.Delete))
}
