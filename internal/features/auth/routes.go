package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *authFeature) PublicRoutes(r chi.Router) {
	r.Post("/auth/login", endpoint.JSON(f.errorHandler, f.authHandler.Login))
}

func (f *authFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/auth/logout", endpoint.JSON(f.errorHandler, f.authHandler.Logout))
}
