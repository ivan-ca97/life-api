package service

import (
	"time"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
	"github.com/ivan-ca97/life/internal/features/ai_usage/ports"
)

type service struct {
	repository ports.Repository
	now        func() time.Time
}

var _ ports.Service = (*service)(nil)

func NewService(repository ports.Repository) *service {
	return &service{repository: repository, now: time.Now}
}

func (s *service) currentPeriod() time.Time {
	return domain.PeriodStart(s.now())
}

func (s *service) CheckQuota(userId uuid.UUID) error {
	allocation, err := s.repository.GetAllocation(userId)
	if err != nil {
		return err
	}
	if !allocation.Tier.Enabled {
		return domain.ErrTierDisabled
	}
	limit := allocation.EffectiveLimitUsd()
	if limit == nil {
		return nil // unlimited
	}
	usage, err := s.repository.GetUsage(userId, s.currentPeriod())
	if err != nil {
		return err
	}
	if domain.OverLimit(usage.CostUsd, limit) {
		return domain.ErrQuotaExceeded
	}
	return nil
}

func (s *service) RecordUsage(userId uuid.UUID, delta ports.UsageDelta) error {
	return s.repository.AddUsage(userId, s.currentPeriod(), delta)
}

func (s *service) GetUsage(userId uuid.UUID) (*domain.UsageSummary, error) {
	allocation, err := s.repository.GetAllocation(userId)
	if err != nil {
		return nil, err
	}
	usage, err := s.repository.GetUsage(userId, s.currentPeriod())
	if err != nil {
		return nil, err
	}
	return &domain.UsageSummary{
		PeriodStart:       usage.PeriodStart,
		Requests:          usage.Requests,
		InputTokens:       usage.InputTokens,
		OutputTokens:      usage.OutputTokens,
		CostUsd:           usage.CostUsd,
		EffectiveLimitUsd: allocation.EffectiveLimitUsd(),
		TierName:          allocation.Tier.Name,
	}, nil
}

func (s *service) SetSelfLimit(userId uuid.UUID, selfLimitUsd *float64) error {
	if selfLimitUsd != nil && *selfLimitUsd < 0 {
		return cerr.NewBadRequestError("self limit cannot be negative")
	}
	return s.repository.SetSelfLimit(userId, selfLimitUsd)
}

func (s *service) ListTiers() ([]domain.Tier, error) {
	return s.repository.ListTiers()
}

func (s *service) CreateTier(params ports.CreateTierParams) (*domain.Tier, error) {
	if params.Name == "" {
		return nil, cerr.NewBadRequestError("tier name is required")
	}
	if params.MonthlyLimitUsd != nil && *params.MonthlyLimitUsd < 0 {
		return nil, cerr.NewBadRequestError("monthly limit cannot be negative")
	}
	tier := &domain.Tier{
		Id:              uuid.New(),
		Name:            params.Name,
		MonthlyLimitUsd: params.MonthlyLimitUsd,
		Enabled:         params.Enabled,
	}
	err := s.repository.CreateTier(tier)
	if err != nil {
		return nil, err
	}
	return tier, nil
}

func (s *service) UpdateTier(id uuid.UUID, params ports.UpdateTierParams) (*domain.Tier, error) {
	if params.Name != nil && *params.Name == "" {
		return nil, cerr.NewBadRequestError("tier name cannot be empty")
	}
	if params.MonthlyLimitUsd != nil && *params.MonthlyLimitUsd != nil && **params.MonthlyLimitUsd < 0 {
		return nil, cerr.NewBadRequestError("monthly limit cannot be negative")
	}
	return s.repository.UpdateTier(id, params)
}

func (s *service) AssignTier(userId, tierId uuid.UUID) error {
	return s.repository.AssignTier(userId, tierId)
}

func (s *service) DeleteTier(id uuid.UUID) error {
	return s.repository.DeleteTier(id)
}

// CostUsd prices the tokens with the rate effective at `at`. Returns 0 when no
// price is on record for the provider/model (treat as "unpriced", not free).
func (s *service) CostUsd(provider, model string, inputTokens, outputTokens int64, at time.Time) (float64, error) {
	price, err := s.repository.FindPrice(provider, model, at)
	if err != nil {
		return 0, err
	}
	if price == nil {
		return 0, nil
	}
	const perMillion = 1_000_000.0
	input := float64(inputTokens) / perMillion * price.InputPerMillion
	output := float64(outputTokens) / perMillion * price.OutputPerMillion
	return input + output, nil
}

func (s *service) LogInteraction(entry ports.InteractionEntry) error {
	entry.InputSummary = truncateRunes(entry.InputSummary, domain.MaxInputSummaryLen)
	return s.repository.InsertInteraction(entry)
}

func (s *service) ListInteractions(filter ports.InteractionFilter) (types.Page[domain.Interaction], error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	return s.repository.ListInteractions(filter)
}

// truncateRunes caps a string to n runes, keeping it valid UTF-8 (Postgres TEXT
// rejects invalid byte sequences, so a plain byte slice could fail on insert).
func truncateRunes(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}
