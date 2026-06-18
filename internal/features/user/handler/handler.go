package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/user/ports"
)

type UserHandler interface {
	Create(r *http.Request) (*userResponse, int, error)
	GetById(r *http.Request) (*userResponse, int, error)
	FindByUsername(r *http.Request) (*userResponse, int, error)
	List(r *http.Request) (*userPage, int, error)
	Update(r *http.Request) (*userResponse, int, error)
	Deactivate(r *http.Request) (*api.NoResponse, int, error)
	AddProfilePhoto(r *http.Request) (*profilePhotoResponse, int, error)
	ListProfilePhotos(r *http.Request) (*profilePhotoPage, int, error)
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
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	user, err := h.service.GetById(r.Context(), userId)
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

func (h *userHandler) FindByUsername(r *http.Request) (*userResponse, int, error) {
	username := r.URL.Query().Get("username")
	if username == "" {
		return nil, 0, cerr.NewBadRequestError("username query parameter is required")
	}
	user, err := h.service.FindByUsername(r.Context(), username)
	if err != nil {
		return nil, 0, err
	}
	return userFromDomain(user), http.StatusOK, nil
}

func (h *userHandler) Update(r *http.Request) (*userResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateUserRequest](r)
	if err != nil {
		return nil, 0, err
	}
	var birthDate *time.Time
	if request.BirthDate != nil {
		t, err := time.Parse("2006-01-02", *request.BirthDate)
		if err != nil {
			return nil, 0, cerr.NewBadRequestError("birth_date must be in YYYY-MM-DD format")
		}
		birthDate = &t
	}
	params := ports.UpdateParams{
		Email:     request.Email,
		Username:  request.Username,
		Password:  request.Password,
		HeightCm:  request.HeightCm,
		BirthDate: birthDate,
		Sex:       request.Sex,
	}
	user, err := h.service.Update(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return userFromDomain(user), http.StatusOK, nil
}

func (h *userHandler) Deactivate(r *http.Request) (*api.NoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	err = h.service.Deactivate(r.Context(), userId)
	if err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}

func (h *userHandler) AddProfilePhoto(r *http.Request) (*profilePhotoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[addProfilePhotoRequest](r)
	if err != nil {
		return nil, 0, err
	}
	photo, err := h.service.AddProfilePhoto(r.Context(), userId, request.Url)
	if err != nil {
		return nil, 0, err
	}
	return profilePhotoFromDomain(photo), http.StatusCreated, nil
}

func (h *userHandler) ListProfilePhotos(r *http.Request) (*profilePhotoPage, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	page, err := h.service.ListProfilePhotos(r.Context(), userId, api.PaginationFromRequest(r))
	if err != nil {
		return nil, 0, err
	}
	return newProfilePhotoPage(page), http.StatusOK, nil
}
