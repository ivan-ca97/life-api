package handler

import "github.com/google/uuid"

type setSelfLimitRequest struct {
	// SelfLimitUsd is the user's optional monthly cap. null/omitted removes it.
	SelfLimitUsd *float64 `json:"self_limit_usd"`
}

type createTierRequest struct {
	Name string `json:"name"`
	// MonthlyLimitUsd null/omitted means unlimited.
	MonthlyLimitUsd *float64 `json:"monthly_limit_usd"`
	Enabled         bool     `json:"enabled"`
}

type updateTierRequest struct {
	Name *string `json:"name,omitempty"`
	// MonthlyLimitUsd sets a numeric cap. To clear it (make unlimited), send
	// "unlimited": true instead.
	MonthlyLimitUsd *float64 `json:"monthly_limit_usd,omitempty"`
	Unlimited       *bool    `json:"unlimited,omitempty"`
	Enabled         *bool    `json:"enabled,omitempty"`
}

type assignTierRequest struct {
	TierId uuid.UUID `json:"tier_id"`
}
