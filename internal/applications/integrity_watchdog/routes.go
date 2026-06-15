package integrity_watchdog

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (a *WatchdogApplication) ProtectedRoutes(r chi.Router) {
	r.Post("/watchdog/trigger", endpoint.JSON(a.errorHandler, a.watchdogHandler.Trigger))
	r.Get("/watchdog/status", endpoint.JSON(a.errorHandler, a.watchdogHandler.Status))
	r.Patch("/watchdog/config", endpoint.JSON(a.errorHandler, a.watchdogHandler.Configure))
}
