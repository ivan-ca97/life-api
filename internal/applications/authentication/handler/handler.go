package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/applications/authentication/ports"
)

type AuthenticationHandler interface {
	Register(r *http.Request) (*authenticationResponse, int, error)
	Login(r *http.Request) (*authenticationResponse, int, error)
	LoginWithGoogle(r *http.Request) (*authenticationResponse, int, error)
	Logout(r *http.Request) (*api.NoResponse, int, error)
}

type authenticationHandler struct {
	useCase ports.AuthenticationUseCase
}

var _ AuthenticationHandler = (*authenticationHandler)(nil)

func NewAuthenticationHandler(useCase ports.AuthenticationUseCase) *authenticationHandler {
	return &authenticationHandler{
		useCase: useCase,
	}
}

func (h *authenticationHandler) Register(r *http.Request) (*authenticationResponse, int, error) {
	request, err := api.DecodeBody[registerRequest](r)
	if err != nil {
		return nil, 0, err
	}

	result, err := h.useCase.Register(request.Email, request.Password)
	if err != nil {
		return nil, 0, err
	}

	response := authenticationResponseFromResult(result)
	return response, http.StatusCreated, nil
}

func (h *authenticationHandler) Login(r *http.Request) (*authenticationResponse, int, error) {
	request, err := api.DecodeBody[loginRequest](r)
	if err != nil {
		return nil, 0, err
	}

	result, err := h.useCase.Login(request.Email, request.Password)
	if err != nil {
		return nil, 0, err
	}

	response := authenticationResponseFromResult(result)
	return response, http.StatusOK, nil
}

func (h *authenticationHandler) LoginWithGoogle(r *http.Request) (*authenticationResponse, int, error) {
	request, err := api.DecodeBody[googleLoginRequest](r)
	if err != nil {
		return nil, 0, err
	}

	result, err := h.useCase.LoginWithGoogle(request.IdToken)
	if err != nil {
		return nil, 0, err
	}

	response := authenticationResponseFromResult(result)
	return response, http.StatusOK, nil
}

func (h *authenticationHandler) Logout(r *http.Request) (*api.NoResponse, int, error) {
	sessionId, err := auth.SessionFromContext(r.Context())
	if err != nil {
		return nil, 0, err
	}

	err = h.useCase.Logout(sessionId)
	if err != nil {
		return nil, 0, err
	}

	return nil, http.StatusNoContent, nil
}
