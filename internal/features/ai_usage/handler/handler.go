package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"

	"github.com/ivan-ca97/life/internal/features/ai_usage/ports"
)

type AiUsageHandler interface {
	GetMyUsage(r *http.Request) (*usageResponse, int, error)
	SetMySelfLimit(r *http.Request) (*usageResponse, int, error)
	ListTiers(r *http.Request) (*tierListResponse, int, error)
	CreateTier(r *http.Request) (*tierResponse, int, error)
	UpdateTier(r *http.Request) (*tierResponse, int, error)
	AssignUserTier(r *http.Request) (*api.NoResponse, int, error)
	GetUserUsage(r *http.Request) (*usageResponse, int, error)
	ListInteractions(r *http.Request) (*interactionListResponse, int, error)
}

type aiUsageHandler struct {
	service ports.AuthorizedService
}

var _ AiUsageHandler = (*aiUsageHandler)(nil)

func NewAiUsageHandler(service ports.AuthorizedService) *aiUsageHandler {
	return &aiUsageHandler{service: service}
}

func (h *aiUsageHandler) GetMyUsage(r *http.Request) (*usageResponse, int, error) {
	summary, err := h.service.GetMyUsage(r.Context())
	if err != nil {
		return nil, 0, err
	}
	return usageFromDomain(summary), http.StatusOK, nil
}

func (h *aiUsageHandler) SetMySelfLimit(r *http.Request) (*usageResponse, int, error) {
	request, err := api.DecodeBody[setSelfLimitRequest](r)
	if err != nil {
		return nil, 0, err
	}
	if err := h.service.SetMySelfLimit(r.Context(), request.SelfLimitUSD); err != nil {
		return nil, 0, err
	}
	summary, err := h.service.GetMyUsage(r.Context())
	if err != nil {
		return nil, 0, err
	}
	return usageFromDomain(summary), http.StatusOK, nil
}

func (h *aiUsageHandler) ListTiers(r *http.Request) (*tierListResponse, int, error) {
	tiers, err := h.service.ListTiers(r.Context())
	if err != nil {
		return nil, 0, err
	}
	items := make([]tierResponse, len(tiers))
	for i, t := range tiers {
		items[i] = tierFromDomain(t)
	}
	return &tierListResponse{Items: items}, http.StatusOK, nil
}

func (h *aiUsageHandler) CreateTier(r *http.Request) (*tierResponse, int, error) {
	request, err := api.DecodeBody[createTierRequest](r)
	if err != nil {
		return nil, 0, err
	}
	tier, err := h.service.CreateTier(r.Context(), ports.CreateTierParams{
		Name:            request.Name,
		MonthlyLimitUSD: request.MonthlyLimitUSD,
		Enabled:         request.Enabled,
	})
	if err != nil {
		return nil, 0, err
	}
	resp := tierFromDomain(*tier)
	return &resp, http.StatusCreated, nil
}

func (h *aiUsageHandler) UpdateTier(r *http.Request) (*tierResponse, int, error) {
	tierId, err := api.PathParamUUID(r, "tierId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateTierRequest](r)
	if err != nil {
		return nil, 0, err
	}
	params := ports.UpdateTierParams{
		Name:    request.Name,
		Enabled: request.Enabled,
	}
	// "unlimited": true clears the cap (NULL); otherwise a numeric value sets it.
	switch {
	case request.Unlimited != nil && *request.Unlimited:
		var none *float64
		params.MonthlyLimitUSD = &none
	case request.MonthlyLimitUSD != nil:
		limit := request.MonthlyLimitUSD
		params.MonthlyLimitUSD = &limit
	}
	tier, err := h.service.UpdateTier(r.Context(), tierId, params)
	if err != nil {
		return nil, 0, err
	}
	resp := tierFromDomain(*tier)
	return &resp, http.StatusOK, nil
}

func (h *aiUsageHandler) AssignUserTier(r *http.Request) (*api.NoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[assignTierRequest](r)
	if err != nil {
		return nil, 0, err
	}
	if err := h.service.AssignUserTier(r.Context(), userId, request.TierId); err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}

func (h *aiUsageHandler) GetUserUsage(r *http.Request) (*usageResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	summary, err := h.service.GetUserUsage(r.Context(), userId)
	if err != nil {
		return nil, 0, err
	}
	return usageFromDomain(summary), http.StatusOK, nil
}

func (h *aiUsageHandler) ListInteractions(r *http.Request) (*interactionListResponse, int, error) {
	userId, err := api.QueryParamUUID(r, "user_id")
	if err != nil {
		return nil, 0, err
	}
	page, err := h.service.ListInteractions(r.Context(), ports.InteractionFilter{
		PaginationParams: api.PaginationFromRequest(r),
		UserId:           userId,
	})
	if err != nil {
		return nil, 0, err
	}
	return interactionsFromPage(page), http.StatusOK, nil
}
