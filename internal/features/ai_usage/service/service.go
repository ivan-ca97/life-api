package service

import (
	"time"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

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
	limit := allocation.EffectiveLimitUSD()
	if limit == nil {
		return nil // unlimited
	}
	usage, err := s.repository.GetUsage(userId, s.currentPeriod())
	if err != nil {
		return err
	}
	if domain.OverLimit(usage.CostUSD, limit) {
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
		CostUSD:           usage.CostUSD,
		EffectiveLimitUSD: allocation.EffectiveLimitUSD(),
		TierName:          allocation.Tier.Name,
	}, nil
}

func (s *service) SetSelfLimit(userId uuid.UUID, selfLimitUSD *float64) error {
	if selfLimitUSD != nil && *selfLimitUSD < 0 {
		return cerr.NewBadRequestError("self limit cannot be negative")
	}
	return s.repository.SetSelfLimit(userId, selfLimitUSD)
}

func (s *service) ListTiers() ([]domain.Tier, error) {
	return s.repository.ListTiers()
}

func (s *service) CreateTier(params ports.CreateTierParams) (*domain.Tier, error) {
	if params.Name == "" {
		return nil, cerr.NewBadRequestError("tier name is required")
	}
	if params.MonthlyLimitUSD != nil && *params.MonthlyLimitUSD < 0 {
		return nil, cerr.NewBadRequestError("monthly limit cannot be negative")
	}
	tier := &domain.Tier{
		Id:              uuid.New(),
		Name:            params.Name,
		MonthlyLimitUSD: params.MonthlyLimitUSD,
		Enabled:         params.Enabled,
	}
	if err := s.repository.CreateTier(tier); err != nil {
		return nil, err
	}
	return tier, nil
}

func (s *service) UpdateTier(id uuid.UUID, params ports.UpdateTierParams) (*domain.Tier, error) {
	if params.Name != nil && *params.Name == "" {
		return nil, cerr.NewBadRequestError("tier name cannot be empty")
	}
	if params.MonthlyLimitUSD != nil && *params.MonthlyLimitUSD != nil && **params.MonthlyLimitUSD < 0 {
		return nil, cerr.NewBadRequestError("monthly limit cannot be negative")
	}
	return s.repository.UpdateTier(id, params)
}

func (s *service) AssignTier(userId, tierId uuid.UUID) error {
	return s.repository.AssignTier(userId, tierId)
}
