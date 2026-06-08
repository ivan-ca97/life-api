package authentication

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (a *authenticationApplication) PublicRoutes(r chi.Router) {
	r.Post("/auth/register", endpoint.JSON(a.errorHandler, a.authenticationHandler.Register))
	r.Post("/auth/login", endpoint.JSON(a.errorHandler, a.authenticationHandler.Login))
	r.Post("/auth/google", endpoint.JSON(a.errorHandler, a.authenticationHandler.LoginWithGoogle))
}

func (a *authenticationApplication) ProtectedRoutes(r chi.Router) {
	r.Post("/auth/logout", endpoint.JSON(a.errorHandler, a.authenticationHandler.Logout))
}
