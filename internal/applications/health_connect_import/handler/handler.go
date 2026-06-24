package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"

	"github.com/ivan-ca97/life/internal/applications/health_connect_import/ports"
)

type HealthConnectImportHandler interface {
	Import(r *http.Request) (*importResponse, int, error)
}

type healthConnectImportHandler struct {
	useCase   ports.HealthConnectImportUseCase
	dumpStore ports.DumpStore // non-nil → dump mode: save payload and short-circuit
}

func NewHealthConnectImportHandler(useCase ports.HealthConnectImportUseCase, dumpStore ports.DumpStore) HealthConnectImportHandler {
	return &healthConnectImportHandler{
		useCase:   useCase,
		dumpStore: dumpStore,
	}
}

func (h *healthConnectImportHandler) Import(r *http.Request) (*importResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}

	payload, err := api.DecodeBody[ports.Payload](r)
	if err != nil {
		return nil, 0, err
	}

	if h.dumpStore != nil {
		raw, err := json.Marshal(payload)
		if err == nil {
			_ = h.dumpStore.Save(userId, payload.AppVersion, raw)
		}
		return &importResponse{}, http.StatusOK, nil
	}

	result, err := h.useCase.Import(r.Context(), userId, payload)
	if err != nil {
		return nil, 0, err
	}

	return importResponseFromResult(result), http.StatusOK, nil
}
