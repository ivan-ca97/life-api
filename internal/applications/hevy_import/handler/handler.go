package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/applications/hevy_import/ports"
)

type HevyImportHandler interface {
	ImportHevy(r *http.Request) (*importResponse, int, error)
}

type hevyImportHandler struct {
	useCase ports.HevyImportUseCase
}

func NewHevyImportHandler(useCase ports.HevyImportUseCase) HevyImportHandler {
	return &hevyImportHandler{
		useCase: useCase,
	}
}

func (h *hevyImportHandler) ImportHevy(r *http.Request) (*importResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, 0, cerr.NewBadRequestError("invalid multipart form")
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, 0, cerr.NewBadRequestError("missing file field")
	}
	defer file.Close()

	result, err := h.useCase.Import(r.Context(), userId, file)
	if err != nil {
		return nil, 0, err
	}

	response := importResponseFromResult(result)
	return response, http.StatusOK, nil
}
