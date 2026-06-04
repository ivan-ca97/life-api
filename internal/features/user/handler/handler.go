package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"

	"github.com/ivan-ca97/life/internal/features/user/ports"
)

type UserHandler interface {
	Create(r *http.Request) (*userResponse, int, error)
	GetById(r *http.Request) (*userResponse, int, error)
	List(r *http.Request) (*userPage, int, error)
	Update(r *http.Request) (*userResponse, int, error)
	Deactivate(r *http.Request) (*api.NoResponse, int, error)
}

type userHandler struct {
	service ports.AuthorizedUserService
}

var _ UserHandler = (*userHandler)(nil)

func NewUserHandler(service ports.AuthorizedUserService) *userHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) Create(r *http.Request) (*userResponse, int, error) {
	request, err := api.DecodeBody[createUserRequest](r)
	if err != nil {
		return nil, 0, err
	}
	user, err := h.service.Create(r.Context(), request.Email, request.Password)
	if err != nil {
		return nil, 0, err
	}
	return userFromDomain(user), http.StatusCreated, nil
}

func (h *userHandler) GetById(r *http.Request) (*userResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	user, err := h.service.GetById(r.Context(), id)
	if err != nil {
		return nil, 0, err
	}
	return userFromDomain(user), http.StatusOK, nil
}

func (h *userHandler) List(r *http.Request) (*userPage, int, error) {
	page, err := h.service.List(r.Context(), api.PaginationFromRequest(r))
	if err != nil {
		return nil, 0, err
	}
	return newUserPage(page), http.StatusOK, nil
}

func (h *userHandler) Update(r *http.Request) (*userResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateUserRequest](r)
	if err != nil {
		return nil, 0, err
	}
	params := ports.UpdateParams{
		Email:    request.Email,
		Password: request.Password,
	}
	user, err := h.service.Update(r.Context(), id, params)
	if err != nil {
		return nil, 0, err
	}
	return userFromDomain(user), http.StatusOK, nil
}

func (h *userHandler) Deactivate(r *http.Request) (*api.NoResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	err = h.service.Deactivate(r.Context(), id)
	if err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}
