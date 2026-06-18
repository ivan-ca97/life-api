package measurements

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *measurementsFeature) PublicRoutes(_ chi.Router) {}

func (f *measurementsFeature) ProtectedRoutes(r chi.Router) {
	r.Put("/body-measurements/{date}/{type}", endpoint.JSON(f.errorHandler, f.handler.Upsert))
	r.Get("/body-measurements/{date}/{type}", endpoint.JSON(f.errorHandler, f.handler.GetByDate))
	r.Get("/body-measurements", endpoint.JSON(f.errorHandler, f.handler.List))
	r.Delete("/body-measurements/{date}/{type}", endpoint.JSON(f.errorHandler, f.handler.Delete))
}
