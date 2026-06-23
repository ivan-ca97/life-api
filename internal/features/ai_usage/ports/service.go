package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
)

type CreateTierParams struct {
	Name            string
	MonthlyLimitUSD *float64
	Enabled         bool
}

// Service is the base (unauthorized) API. The meal AI use case depends on the
// QuotaGuard subset; HTTP handlers go through AuthorizedService.
type Service interface {
	QuotaGuard

	GetUsage(userId uuid.UUID) (*domain.UsageSummary, error)
	SetSelfLimit(userId uuid.UUID, selfLimitUSD *float64) error

	ListTiers() ([]domain.Tier, error)
	CreateTier(params CreateTierParams) (*domain.Tier, error)
	UpdateTier(id uuid.UUID, params UpdateTierParams) (*domain.Tier, error)
	AssignTier(userId, tierId uuid.UUID) error
}

// QuotaGuard is the narrow contract the meal AI feature consumes: check before
// spending, record after. Defined here so meal_ai can depend on a small port.
type QuotaGuard interface {
	// CheckQuota returns domain.ErrQuotaExceeded if the user has reached their
	// effective monthly limit.
	CheckQuota(userId uuid.UUID) error
	// RecordUsage adds one call's consumption to the current period.
	RecordUsage(userId uuid.UUID, delta UsageDelta) error
}

// AuthorizedService wraps Service with access control: "me" operations act on
// the actor from context; tier administration requires the admin role.
type AuthorizedService interface {
	GetMyUsage(ctx context.Context) (*domain.UsageSummary, error)
	SetMySelfLimit(ctx context.Context, selfLimitUSD *float64) error

	ListTiers(ctx context.Context) ([]domain.Tier, error)
	CreateTier(ctx context.Context, params CreateTierParams) (*domain.Tier, error)
	UpdateTier(ctx context.Context, id uuid.UUID, params UpdateTierParams) (*domain.Tier, error)
	AssignUserTier(ctx context.Context, userId, tierId uuid.UUID) error
	GetUserUsage(ctx context.Context, userId uuid.UUID) (*domain.UsageSummary, error)
}
