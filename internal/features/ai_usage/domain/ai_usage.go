package domain

import (
	"time"

	"github.com/google/uuid"
)

// Tier is an AI spend plan managed by an admin. MonthlyLimitUSD is nil for an
// unlimited tier.
type Tier struct {
	Id              uuid.UUID
	Name            string
	MonthlyLimitUSD *float64
	IsDefault       bool
	Enabled         bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Allocation is a user's effective plan: their tier plus an optional
// self-imposed cap they can set for themselves.
type Allocation struct {
	Tier         Tier
	SelfLimitUSD *float64
}

// EffectiveLimitUSD is the binding monthly spend cap in USD, or nil when the
// user is effectively unlimited. It is the smaller of the tier limit and the
// user's self-imposed limit, ignoring whichever is unset. This is what lets an
// "unlimited" user still protect themselves with a self-limit.
func (a Allocation) EffectiveLimitUSD() *float64 {
	return minLimit(a.Tier.MonthlyLimitUSD, a.SelfLimitUSD)
}

func minLimit(a, b *float64) *float64 {
	switch {
	case a == nil:
		return b
	case b == nil:
		return a
	case *a <= *b:
		return a
	default:
		return b
	}
}

// Usage is a user's accumulated AI consumption within one monthly period.
type Usage struct {
	UserId       uuid.UUID
	PeriodStart  time.Time
	Requests     int
	InputTokens  int64
	OutputTokens int64
	CostUSD      float64
}

// UsageSummary is what the user (or an admin) sees: consumption plus the
// effective limit it is measured against.
type UsageSummary struct {
	PeriodStart       time.Time
	Requests          int
	InputTokens       int64
	OutputTokens      int64
	CostUSD           float64
	EffectiveLimitUSD *float64
	TierName          string
}

// OverLimit reports whether an accumulated cost has reached a limit. A nil limit
// is unlimited and never over.
func OverLimit(costUSD float64, limit *float64) bool {
	return limit != nil && costUSD >= *limit
}

// PeriodStart returns the first day (UTC, date-only) of the month containing t.
// It is the partition key for usage accounting.
func PeriodStart(t time.Time) time.Time {
	y, m, _ := t.UTC().Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}
