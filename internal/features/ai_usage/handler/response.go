package handler

import (
	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
)

type tierResponse struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	MonthlyLimitUSD *float64 `json:"monthly_limit_usd"`
	IsDefault       bool     `json:"is_default"`
	Enabled         bool     `json:"enabled"`
}

func tierFromDomain(t domain.Tier) tierResponse {
	return tierResponse{
		Id:              t.Id.String(),
		Name:            t.Name,
		MonthlyLimitUSD: t.MonthlyLimitUSD,
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
	CostUSD           float64  `json:"cost_usd"`
	EffectiveLimitUSD *float64 `json:"effective_limit_usd"`
	TierName          string   `json:"tier_name"`
}

func usageFromDomain(u *domain.UsageSummary) *usageResponse {
	return &usageResponse{
		PeriodStart:       u.PeriodStart.Format("2006-01-02"),
		Requests:          u.Requests,
		InputTokens:       u.InputTokens,
		OutputTokens:      u.OutputTokens,
		CostUSD:           u.CostUSD,
		EffectiveLimitUSD: u.EffectiveLimitUSD,
		TierName:          u.TierName,
	}
}
