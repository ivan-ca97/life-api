package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	"github.com/ivan-ca97/life/pkg/auth"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/scheduler"
)

type WatchdogHandler interface {
	Trigger(r *http.Request) (*statusResponse, int, error)
	Status(r *http.Request) (*statusResponse, int, error)
	Configure(r *http.Request) (*statusResponse, int, error)
}

type watchdogHandler struct {
	scheduler  *scheduler.Scheduler
	authorizer auth.AuthorizationService
}

func NewWatchdogHandler(s *scheduler.Scheduler, authorizer auth.AuthorizationService) WatchdogHandler {
	return &watchdogHandler{scheduler: s, authorizer: authorizer}
}

func (h *watchdogHandler) Trigger(r *http.Request) (*statusResponse, int, error) {
	err := h.authorizer.AuthorizeAdmin(r.Context())
	if err != nil {
		return nil, 0, err
	}
	if !h.scheduler.Trigger() {
		return nil, 0, cerr.NewConflictError("a run is already pending or in progress")
	}
	return buildStatusResponse(h.scheduler), http.StatusAccepted, nil
}

func (h *watchdogHandler) Status(r *http.Request) (*statusResponse, int, error) {
	err := h.authorizer.AuthorizeAdmin(r.Context())
	if err != nil {
		return nil, 0, err
	}
	return buildStatusResponse(h.scheduler), http.StatusOK, nil
}

func (h *watchdogHandler) Configure(r *http.Request) (*statusResponse, int, error) {
	err := h.authorizer.AuthorizeAdmin(r.Context())
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[configureRequest](r)
	if err != nil {
		return nil, 0, err
	}
	if request.IntervalSeconds <= 0 {
		return nil, 0, cerr.NewBadRequestError("interval_seconds must be positive")
	}
	h.scheduler.SetPeriod(time.Duration(request.IntervalSeconds) * time.Second)
	return buildStatusResponse(h.scheduler), http.StatusOK, nil
}
