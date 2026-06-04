package daily

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *dailyFeature) PublicRoutes(_ chi.Router) {}

func (f *dailyFeature) ProtectedRoutes(r chi.Router) {
	r.Get("/daily/summary", endpoint.JSON(f.errorHandler, f.summaryHandler.GetSummary))
	r.Get("/daily/summary/range", endpoint.JSON(f.errorHandler, f.summaryHandler.GetSummaryRange))
	r.Get("/daily/corrections", endpoint.JSON(f.errorHandler, f.correctionHandler.GetCorrection))
	r.Put("/daily/corrections", endpoint.JSON(f.errorHandler, f.correctionHandler.UpsertCorrection))
	r.Delete("/daily/corrections", endpoint.JSON(f.errorHandler, f.correctionHandler.DeleteCorrection))
}
