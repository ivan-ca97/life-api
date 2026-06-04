package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/auth/ports"
)

type AuthHandler interface {
	Login(r *http.Request) (*loginResponse, int, error)
	Logout(r *http.Request) (*api.NoResponse, int, error)
}

type authHandler struct {
	service ports.AuthService
}

var _ AuthHandler = (*authHandler)(nil)

func NewAuthHandler(service ports.AuthService) *authHandler {
	return &authHandler{
		service: service,
	}
}

func (h *authHandler) Login(r *http.Request) (*loginResponse, int, error) {
	request, err := api.DecodeBody[loginRequest](r)
	if err != nil {
		return nil, 0, err
	}
	session, err := h.service.Login(request.Email, request.Password)
	if err != nil {
		return nil, 0, err
	}
	response := &loginResponse{
		Token:     session.Id,
		ExpiresAt: session.ExpiresAt,
	}

	return response, http.StatusOK, nil
}

func (h *authHandler) Logout(r *http.Request) (*api.NoResponse, int, error) {
	sessionId, err := auth.SessionFromContext(r.Context())
	if err != nil {
		return nil, 0, err
	}
	err = h.service.Logout(sessionId)
	if err != nil {
		return nil, 0, err
	}

	return nil, http.StatusNoContent, nil
}
