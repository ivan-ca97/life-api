package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type PhotoHandler interface {
	Create(r *http.Request) (*dailyPhotoResponse, int, error)
	List(r *http.Request) (*dailyPhotoListResponse, int, error)
	Update(r *http.Request) (*dailyPhotoResponse, int, error)
	Delete(r *http.Request) (*api.NoResponse, int, error)
}

type photoHandler struct {
	service ports.AuthorizedPhotoService
}

var _ PhotoHandler = (*photoHandler)(nil)

func NewPhotoHandler(service ports.AuthorizedPhotoService) *photoHandler {
	return &photoHandler{service: service}
}

type createDailyPhotoRequest struct {
	Date      string `json:"date"`
	Url       string `json:"url"`
	Name      string `json:"name"`
	IsPrimary bool   `json:"is_primary"`
}

type updateDailyPhotoRequest struct {
	Name      *string `json:"name,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
}

type dailyPhotoResponse struct {
	Id        uuid.UUID `json:"id"`
	Date      string    `json:"date"`
	Url       string    `json:"url"`
	Name      string    `json:"name"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

type dailyPhotoListResponse struct {
	Items []dailyPhotoResponse `json:"items"`
}

func photoFromDomain(p *domain.DailyPhoto) *dailyPhotoResponse {
	return &dailyPhotoResponse{
		Id:        p.Id,
		Date:      p.Date.Format("2006-01-02"),
		Url:       p.Url,
		Name:      p.Name,
		IsPrimary: p.IsPrimary,
		CreatedAt: p.CreatedAt,
	}
}

func (h *photoHandler) Create(r *http.Request) (*dailyPhotoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[createDailyPhotoRequest](r)
	if err != nil {
		return nil, 0, err
	}
	if request.Url == "" {
		return nil, 0, cerr.NewBadRequestError("url is required")
	}
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
	}
	photo, err := h.service.Create(r.Context(), userId, ports.CreatePhotoParams{
		Date:      date,
		Url:       request.Url,
		Name:      request.Name,
		IsPrimary: request.IsPrimary,
	})
	if err != nil {
		return nil, 0, err
	}
	return photoFromDomain(photo), http.StatusCreated, nil
}

func (h *photoHandler) List(r *http.Request) (*dailyPhotoListResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, err := api.QueryParamDate(r, "date")
	if err != nil {
		return nil, 0, err
	}
	if date == nil {
		return nil, 0, cerr.NewBadRequestError("date query parameter is required (format: YYYY-MM-DD)")
	}
	photos, err := h.service.List(r.Context(), userId, *date)
	if err != nil {
		return nil, 0, err
	}
	items := make([]dailyPhotoResponse, len(photos))
	for i, p := range photos {
		items[i] = *photoFromDomain(&p)
	}
	result := &dailyPhotoListResponse{
		Items: items,
	}
	return result, http.StatusOK, nil
}

func (h *photoHandler) Update(r *http.Request) (*dailyPhotoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateDailyPhotoRequest](r)
	if err != nil {
		return nil, 0, err
	}
	photo, err := h.service.Update(r.Context(), userId, id, ports.UpdatePhotoParams{
		Name:      request.Name,
		IsPrimary: request.IsPrimary,
	})
	if err != nil {
		return nil, 0, err
	}
	return photoFromDomain(photo), http.StatusOK, nil
}

func (h *photoHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
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
