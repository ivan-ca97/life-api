package goal

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *goalFeature) PublicRoutes(_ chi.Router) {}

func (f *goalFeature) ProtectedRoutes(r chi.Router) {
	r.Get("/goals", endpoint.JSON(f.errorHandler, f.goalHandler.GetCurrent))
	r.Put("/goals", endpoint.JSON(f.errorHandler, f.goalHandler.Upsert))
}
