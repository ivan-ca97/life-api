package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *authFeature) PublicRoutes(r chi.Router) {
	r.Post("/auth/register", endpoint.JSON(f.errorHandler, f.authHandler.Register))
	r.Post("/auth/login", endpoint.JSON(f.errorHandler, f.authHandler.Login))
	r.Post("/auth/google", endpoint.JSON(f.errorHandler, f.authHandler.LoginWithGoogle))
}

func (f *authFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/auth/logout", endpoint.JSON(f.errorHandler, f.authHandler.Logout))
}
