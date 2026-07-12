package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"

	"github.com/ivan-ca97/life/internal/features/media/ports"
)

type MediaHandler interface {
	GenerateUploadURL(r *http.Request) (*uploadURLResponse, int, error)
}

type mediaHandler struct {
	service ports.MediaService
}

func NewMediaHandler(service ports.MediaService) *mediaHandler {
	return &mediaHandler{
		service: service,
	}
}

func (h *mediaHandler) GenerateUploadURL(r *http.Request) (*uploadURLResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}

	request, err := api.DecodeBody[uploadURLRequest](r)
	if err != nil {
		return nil, 0, err
	}

	result, err := h.service.GenerateUploadURL(r.Context(), ports.UploadRequest{
		UserId:      userId,
		Filename:    request.Filename,
		ContentType: request.ContentType,
	})
	if err != nil {
		return nil, 0, err
	}

	response := &uploadURLResponse{
		UploadURL: result.UploadURL,
		PublicURL: result.PublicURL,
	}
	return response, http.StatusOK, nil
}
