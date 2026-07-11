package handler

import (
	"encoding/json"
	"time"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
)

type tierResponse struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	MonthlyLimitUsd *float64 `json:"monthly_limit_usd"`
	IsDefault       bool     `json:"is_default"`
	Enabled         bool     `json:"enabled"`
}

func tierFromDomain(t domain.Tier) tierResponse {
	return tierResponse{
		Id:              t.Id.String(),
		Name:            t.Name,
		MonthlyLimitUsd: t.MonthlyLimitUsd,
		IsDefault:       t.IsDefault,
		Enabled:         t.Enabled,
	}
}

type tierListResponse struct {
	Items []tierResponse `json:"items"`
}

type usageResponse struct {
	PeriodStart       string   `json:"period_start"`
	Requests          int      `json:"requests"`
	InputTokens       int64    `json:"input_tokens"`
	OutputTokens      int64    `json:"output_tokens"`
	CostUsd           float64  `json:"cost_usd"`
	EffectiveLimitUsd *float64 `json:"effective_limit_usd"`
	TierName          string   `json:"tier_name"`
}

func usageFromDomain(u *domain.UsageSummary) *usageResponse {
	return &usageResponse{
		PeriodStart:       u.PeriodStart.Format("2006-01-02"),
		Requests:          u.Requests,
		InputTokens:       u.InputTokens,
		OutputTokens:      u.OutputTokens,
		CostUsd:           u.CostUsd,
		EffectiveLimitUsd: u.EffectiveLimitUsd,
		TierName:          u.TierName,
	}
}

type interactionResponse struct {
	Id            string          `json:"id"`
	UserId        string          `json:"user_id"`
	CreatedAt     time.Time       `json:"created_at"`
	Operation     string          `json:"operation"`
	Provider      string          `json:"provider"`
	Model         string          `json:"model"`
	Status        string          `json:"status"`
	ErrorType     string          `json:"error_type,omitempty"`
	InputTokens   int64           `json:"input_tokens"`
	OutputTokens  int64           `json:"output_tokens"`
	CostUsd       float64         `json:"cost_usd"`
	LatencyMs     int             `json:"latency_ms"`
	ProviderCalls int             `json:"provider_calls"`
	CorrelationId *string         `json:"correlation_id,omitempty"`
	InputSummary  string          `json:"input_summary,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
}

type interactionListResponse struct {
	Items  []interactionResponse `json:"items"`
	Total  int64                 `json:"total"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}

func interactionsFromPage(page types.Page[domain.Interaction]) *interactionListResponse {
	items := make([]interactionResponse, len(page.Items))
	for i, it := range page.Items {
		var correlationId *string
		if it.CorrelationId != nil {
			s := it.CorrelationId.String()
			correlationId = &s
		}
		items[i] = interactionResponse{
			Id:            it.Id.String(),
			UserId:        it.UserId.String(),
			CreatedAt:     it.CreatedAt,
			Operation:     it.Operation,
			Provider:      it.Provider,
			Model:         it.Model,
			Status:        it.Status,
			ErrorType:     it.ErrorType,
			InputTokens:   it.InputTokens,
			OutputTokens:  it.OutputTokens,
			CostUsd:       it.CostUsd,
			LatencyMs:     it.LatencyMs,
			ProviderCalls: it.ProviderCalls,
			CorrelationId: correlationId,
			InputSummary:  it.InputSummary,
			Metadata:      it.Metadata,
		}
	}
	return &interactionListResponse{
		Items:  items,
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
}
