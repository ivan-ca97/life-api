package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"

	"github.com/ivan-ca97/life/internal/applications/authorization/ports"
)

type ShareHandler interface {
	Create(r *http.Request) (*shareResponse, int, error)
	ListOwned(r *http.Request) (*shareListResponse, int, error)
	ListReceived(r *http.Request) (*shareListResponse, int, error)
	Update(r *http.Request) (*shareResponse, int, error)
	Delete(r *http.Request) (*api.NoResponse, int, error)
}

type shareHandler struct {
	service ports.AuthorizedShareUseCase
}

var _ ShareHandler = (*shareHandler)(nil)

func NewShareHandler(service ports.AuthorizedShareUseCase) *shareHandler {
	return &shareHandler{
		service: service,
	}
}

func (h *shareHandler) Create(r *http.Request) (*shareResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[createShareRequest](r)
	if err != nil {
		return nil, 0, err
	}
	share, err := h.service.Create(r.Context(), userId, request.GranteeEmail, request.ResourceType, request.CanWrite)
	if err != nil {
		return nil, 0, err
	}
	return shareFromDomain(share), http.StatusCreated, nil
}

func (h *shareHandler) ListOwned(r *http.Request) (*shareListResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	shares, err := h.service.ListOwned(r.Context(), userId)
	if err != nil {
		return nil, 0, err
	}
	return shareListFromDomain(shares), http.StatusOK, nil
}

func (h *shareHandler) ListReceived(r *http.Request) (*shareListResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	shares, err := h.service.ListReceived(r.Context(), userId)
	if err != nil {
		return nil, 0, err
	}
	return shareListFromDomain(shares), http.StatusOK, nil
}

func (h *shareHandler) Update(r *http.Request) (*shareResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateShareRequest](r)
	if err != nil {
		return nil, 0, err
	}
	share, err := h.service.Update(r.Context(), userId, id, request.CanWrite)
	if err != nil {
		return nil, 0, err
	}
	return shareFromDomain(share), http.StatusOK, nil
}

func (h *shareHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	err = h.service.Delete(r.Context(), userId, id)
	if err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}
