package handler

import "github.com/google/uuid"

type setSelfLimitRequest struct {
	// SelfLimitUSD is the user's optional monthly cap. null/omitted removes it.
	SelfLimitUSD *float64 `json:"self_limit_usd"`
}

type createTierRequest struct {
	Name string `json:"name"`
	// MonthlyLimitUSD null/omitted means unlimited.
	MonthlyLimitUSD *float64 `json:"monthly_limit_usd"`
	Enabled         bool     `json:"enabled"`
}

type updateTierRequest struct {
	Name *string `json:"name,omitempty"`
	// MonthlyLimitUSD sets a numeric cap. To clear it (make unlimited), send
	// "unlimited": true instead.
	MonthlyLimitUSD *float64 `json:"monthly_limit_usd,omitempty"`
	Unlimited       *bool    `json:"unlimited,omitempty"`
	Enabled         *bool    `json:"enabled,omitempty"`
}

type assignTierRequest struct {
	TierId uuid.UUID `json:"tier_id"`
}
