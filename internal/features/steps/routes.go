package steps

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *stepsFeature) ProtectedRoutes(r chi.Router) {
	r.Put("/steps/{date}", endpoint.JSON(f.errorHandler, f.handler.Upsert))
	r.Get("/steps", endpoint.JSON(f.errorHandler, f.handler.List))
	r.Get("/steps/{date}", endpoint.JSON(f.errorHandler, f.handler.GetByDate))
	r.Delete("/steps/{date}", endpoint.JSON(f.errorHandler, f.handler.Delete))
}
