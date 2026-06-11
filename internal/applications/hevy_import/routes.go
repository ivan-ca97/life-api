package hevy_import

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (a *hevyImportApplication) ProtectedRoutes(r chi.Router) {
	r.Post("/exercises/import/hevy", endpoint.JSON(a.errorHandler, a.importHandler.ImportHevy))
}
