package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"

	"github.com/ivan-ca97/life/internal/features/media/ports"
)

type MediaHandler interface {
	GenerateUploadUrl(r *http.Request) (*uploadUrlResponse, int, error)
}

type mediaHandler struct {
	service ports.AuthorizedMediaService
}

var _ MediaHandler = (*mediaHandler)(nil)

func NewMediaHandler(service ports.AuthorizedMediaService) *mediaHandler {
	return &mediaHandler{
		service: service,
	}
}

func (h *mediaHandler) GenerateUploadUrl(r *http.Request) (*uploadUrlResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}

	request, err := api.DecodeBody[uploadUrlRequest](r)
	if err != nil {
		return nil, 0, err
	}

	uploadRequest := ports.UploadRequest{
		UserId:      userId,
		Filename:    request.Filename,
		ContentType: request.ContentType,
	}
	result, err := h.service.GenerateUploadUrl(r.Context(), uploadRequest)
	if err != nil {
		return nil, 0, err
	}

	response := &uploadUrlResponse{
		UploadUrl: result.UploadUrl,
		PublicUrl: result.PublicUrl,
	}
	return response, http.StatusOK, nil
}
