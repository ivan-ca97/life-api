package authorization

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (a *authorizationApplication) ProtectedRoutes(r chi.Router) {
	r.Post("/shares", endpoint.JSON(a.errorHandler, a.shareHandler.Create))
	r.Get("/shares", endpoint.JSON(a.errorHandler, a.shareHandler.ListOwned))
	r.Get("/shares/received", endpoint.JSON(a.errorHandler, a.shareHandler.ListReceived))
	r.Patch("/shares/{id}", endpoint.JSON(a.errorHandler, a.shareHandler.Update))
	r.Delete("/shares/{id}", endpoint.JSON(a.errorHandler, a.shareHandler.Delete))
}
