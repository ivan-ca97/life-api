package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"
)

func DecodeBody[T any](r *http.Request) (*T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, cerr.NewBadRequestError("invalid request body")
	}
	return &v, nil
}

func PathParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func PathParamUUID(r *http.Request, key string) (uuid.UUID, error) {
	raw := chi.URLParam(r, key)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.UUID{}, cerr.NewBadRequestError(fmt.Sprintf("invalid %s", key))
	}
	return id, nil
}

func QueryParamDate(r *http.Request, key string) (*time.Time, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}
	date, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, cerr.NewBadRequestError(fmt.Sprintf("invalid %s format, expected YYYY-MM-DD", key))
	}
	return &date, nil
}

func QueryParamInt(r *http.Request, key string) (*int, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return nil, cerr.NewBadRequestError(fmt.Sprintf("invalid %s, must be an integer", key))
	}
	return &v, nil
}

func PaginationFromRequest(r *http.Request) types.PaginationParams {
	params := types.PaginationParams{
		Limit:  20,
		Offset: 0,
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			params.Limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			params.Offset = v
		}
	}
	return params
}
