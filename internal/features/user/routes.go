package user

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *userFeature) PublicRoutes(_ chi.Router) {}

func (f *userFeature) AdminRoutes(r chi.Router) {
	r.Post("/users", endpoint.JSON(f.errorHandler, f.userHandler.Create))
	r.Get("/users", endpoint.JSON(f.errorHandler, f.userHandler.List))
}

func (f *userFeature) ProtectedRoutes(r chi.Router) {
	r.Get("/", endpoint.JSON(f.errorHandler, f.userHandler.GetById))
	r.Patch("/", endpoint.JSON(f.errorHandler, f.userHandler.Update))
	r.Delete("/", endpoint.JSON(f.errorHandler, f.userHandler.Deactivate))
}
