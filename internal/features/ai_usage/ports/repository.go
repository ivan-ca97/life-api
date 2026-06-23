package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
)

type UpdateTierParams struct {
	Name            *string
	MonthlyLimitUSD **float64 // double pointer: distinguishes "unset" from "set to NULL (unlimited)"
	Enabled         *bool
}

// UsageDelta is one increment of consumption to record after an AI call.
type UsageDelta struct {
	Requests     int
	InputTokens  int64
	OutputTokens int64
	CostUSD      float64
}

type Repository interface {
	// Tiers
	ListTiers() ([]domain.Tier, error)
	CreateTier(tier *domain.Tier) error
	UpdateTier(id uuid.UUID, params UpdateTierParams) (*domain.Tier, error)
	GetDefaultTier() (*domain.Tier, error)

	// Allocation (tier + self limit) per user. Falls back to the default tier
	// when the user has no explicit assignment.
	GetAllocation(userId uuid.UUID) (*domain.Allocation, error)
	AssignTier(userId, tierId uuid.UUID) error
	SetSelfLimit(userId uuid.UUID, selfLimitUSD *float64) error

	// Usage accounting, partitioned by month.
	GetUsage(userId uuid.UUID, periodStart time.Time) (*domain.Usage, error)
	AddUsage(userId uuid.UUID, periodStart time.Time, delta UsageDelta) error
}
