package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
)

type CreateTierParams struct {
	Name            string
	MonthlyLimitUsd *float64
	Enabled         bool
}

// Service is the base (unauthorized) API. The meal AI use case depends on the
// QuotaGuard and InteractionLogger subsets; HTTP handlers go through
// AuthorizedService.
type Service interface {
	QuotaGuard
	InteractionLogger

	GetUsage(userId uuid.UUID) (*domain.UsageSummary, error)
	SetSelfLimit(userId uuid.UUID, selfLimitUsd *float64) error

	ListTiers() ([]domain.Tier, error)
	CreateTier(params CreateTierParams) (*domain.Tier, error)
	UpdateTier(id uuid.UUID, params UpdateTierParams) (*domain.Tier, error)
	DeleteTier(id uuid.UUID) error
	AssignTier(userId, tierId uuid.UUID) error

	ListInteractions(filter InteractionFilter) (types.Page[domain.Interaction], error)

	// CostUsd prices token usage with the rate effective at the given time.
	CostUsd(provider, model string, inputTokens, outputTokens int64, at time.Time) (float64, error)
}

// QuotaGuard is the narrow contract the meal AI feature consumes to enforce
// spend limits: check before spending, record after.
type QuotaGuard interface {
	// CheckQuota returns domain.ErrQuotaExceeded if the user has reached their
	// effective monthly limit.
	CheckQuota(userId uuid.UUID) error
	// RecordUsage adds one call's consumption to the current period.
	RecordUsage(userId uuid.UUID, delta UsageDelta) error
}

// InteractionLogger is the narrow contract the meal AI feature consumes to
// record one interaction (best-effort; a failure must not fail the operation).
type InteractionLogger interface {
	LogInteraction(entry InteractionEntry) error
}

// AuthorizedService wraps Service with access control: "me" operations act on
// the actor from context; tier administration requires the admin role.
type AuthorizedService interface {
	GetMyUsage(ctx context.Context) (*domain.UsageSummary, error)
	SetMySelfLimit(ctx context.Context, selfLimitUsd *float64) error

	ListTiers(ctx context.Context) ([]domain.Tier, error)
	CreateTier(ctx context.Context, params CreateTierParams) (*domain.Tier, error)
	UpdateTier(ctx context.Context, id uuid.UUID, params UpdateTierParams) (*domain.Tier, error)
	DeleteTier(ctx context.Context, id uuid.UUID) error
	AssignUserTier(ctx context.Context, userId, tierId uuid.UUID) error
	GetUserUsage(ctx context.Context, userId uuid.UUID) (*domain.UsageSummary, error)
	ListInteractions(ctx context.Context, filter InteractionFilter) (types.Page[domain.Interaction], error)
}
