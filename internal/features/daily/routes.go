package daily

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *dailyFeature) PublicRoutes(_ chi.Router) {}

func (f *dailyFeature) ProtectedRoutes(r chi.Router) {
	r.Get("/daily/summary", endpoint.JSON(f.errorHandler, f.summaryHandler.GetSummary))
	r.Get("/daily/summary/range", endpoint.JSON(f.errorHandler, f.summaryHandler.GetSummaryRange))
	r.Get("/daily/check", endpoint.JSON(f.errorHandler, f.summaryHandler.GetDailyCheck))
	r.Put("/daily/closure", endpoint.JSON(f.errorHandler, f.closureHandler.Close))
	r.Delete("/daily/closure", endpoint.JSON(f.errorHandler, f.closureHandler.Open))
	r.Get("/daily/closure", endpoint.JSON(f.errorHandler, f.closureHandler.GetStatus))
	r.Post("/daily/photos", endpoint.JSON(f.errorHandler, f.photoHandler.Create))
	r.Get("/daily/photos", endpoint.JSON(f.errorHandler, f.photoHandler.List))
	r.Patch("/daily/photos/{id}", endpoint.JSON(f.errorHandler, f.photoHandler.Update))
	r.Delete("/daily/photos/{id}", endpoint.JSON(f.errorHandler, f.photoHandler.Delete))
}
