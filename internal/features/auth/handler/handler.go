package handler

import (
	"errors"
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/auth/ports"
	user_domain "github.com/ivan-ca97/life/internal/features/user/domain"
	user_ports "github.com/ivan-ca97/life/internal/features/user/ports"
)

type AuthHandler interface {
	Register(r *http.Request) (*registerResponse, int, error)
	Login(r *http.Request) (*loginResponse, int, error)
	LoginWithGoogle(r *http.Request) (*loginResponse, int, error)
	Logout(r *http.Request) (*api.NoResponse, int, error)
}

type authHandler struct {
	service        ports.AuthService
	userService    user_ports.UserService
	roleAssigner   ports.RoleAssigner
	googleVerifier ports.GoogleTokenVerifier
	googleClientId string
}

var _ AuthHandler = (*authHandler)(nil)

func NewAuthHandler(
	service ports.AuthService,
	userService user_ports.UserService,
	roleAssigner ports.RoleAssigner,
	googleVerifier ports.GoogleTokenVerifier,
	googleClientId string,
) *authHandler {
	return &authHandler{
		service:        service,
		userService:    userService,
		roleAssigner:   roleAssigner,
		googleVerifier: googleVerifier,
		googleClientId: googleClientId,
	}
}

func (h *authHandler) Register(r *http.Request) (*registerResponse, int, error) {
	request, err := api.DecodeBody[registerRequest](r)
	if err != nil {
		return nil, 0, err
	}
	user, err := h.userService.Create(request.Email, request.Password)
	if err != nil {
		return nil, 0, err
	}
	err = h.roleAssigner.AssignRoleByName(user.Id, "user")
	if err != nil {
		return nil, 0, err
	}
	session, err := h.service.CreateSession(user.Id)
	if err != nil {
		return nil, 0, err
	}
	response := &registerResponse{
		UserId:    user.Id,
		Token:     session.Id,
		ExpiresAt: session.ExpiresAt,
	}
	return response, http.StatusCreated, nil
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
		UserId:    session.UserId,
		Token:     session.Id,
		ExpiresAt: session.ExpiresAt,
	}

	return response, http.StatusOK, nil
}

func (h *authHandler) LoginWithGoogle(r *http.Request) (*loginResponse, int, error) {
	request, err := api.DecodeBody[googleLoginRequest](r)
	if err != nil {
		return nil, 0, err
	}

	claims, err := h.googleVerifier.Verify(request.IdToken, h.googleClientId)
	if err != nil {
		return nil, 0, err
	}

	user, err := h.userService.GetByEmail(claims.Email)
	if err != nil {
		if !errors.Is(err, user_domain.ErrUserNotFound) {
			return nil, 0, err
		}

		user, err = h.userService.CreateOAuth(claims.Email, claims.Subject)
		if err != nil {
			return nil, 0, err
		}
		err = h.roleAssigner.AssignRoleByName(user.Id, "user")
		if err != nil {
			return nil, 0, err
		}
	}

	if !user.Active {
		return nil, 0, user_domain.ErrUserInactive
	}

	session, err := h.service.CreateSession(user.Id)
	if err != nil {
		return nil, 0, err
	}

	response := &loginResponse{
		UserId:    session.UserId,
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
