package health_connect_import

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (a *healthConnectImportApplication) ProtectedRoutes(r chi.Router) {
	r.Post("/import/health-connect", endpoint.JSON(a.errorHandler, a.importHandler.Import))
}
