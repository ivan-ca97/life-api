package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
	"github.com/ivan-ca97/life/internal/features/ai_usage/ports"
)

type authorizedService struct {
	base       ports.Service
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedService = (*authorizedService)(nil)

func NewAuthorizedService(base ports.Service, authorizer auth.AuthorizationService) *authorizedService {
	return &authorizedService{
		base:       base,
		authorizer: authorizer,
	}
}

// --- "me" operations: act on the authenticated actor ---

func (s *authorizedService) GetMyUsage(ctx context.Context) (*domain.UsageSummary, error) {
	actorId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	usage, err := s.base.GetUsage(actorId)
	if err != nil {
		return nil, err
	}
	return usage, nil
}

func (s *authorizedService) SetMySelfLimit(ctx context.Context, selfLimitUsd *float64) error {
	actorId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return err
	}
	err = s.base.SetSelfLimit(actorId, selfLimitUsd)
	if err != nil {
		return err
	}
	return nil
}

// --- admin operations: require the admin role ---

func (s *authorizedService) ListTiers(ctx context.Context) ([]domain.Tier, error) {
	err := s.authorizer.AuthorizeAdmin(ctx)
	if err != nil {
		return nil, err
	}
	tiers, err := s.base.ListTiers()
	if err != nil {
		return nil, err
	}
	return tiers, nil
}

func (s *authorizedService) CreateTier(ctx context.Context, params ports.CreateTierParams) (*domain.Tier, error) {
	err := s.authorizer.AuthorizeAdmin(ctx)
	if err != nil {
		return nil, err
	}
	tier, err := s.base.CreateTier(params)
	if err != nil {
		return nil, err
	}
	return tier, nil
}

func (s *authorizedService) UpdateTier(ctx context.Context, id uuid.UUID, params ports.UpdateTierParams) (*domain.Tier, error) {
	err := s.authorizer.AuthorizeAdmin(ctx)
	if err != nil {
		return nil, err
	}
	tier, err := s.base.UpdateTier(id, params)
	if err != nil {
		return nil, err
	}
	return tier, nil
}

func (s *authorizedService) DeleteTier(ctx context.Context, id uuid.UUID) error {
	err := s.authorizer.AuthorizeAdmin(ctx)
	if err != nil {
		return err
	}
	err = s.base.DeleteTier(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *authorizedService) AssignUserTier(ctx context.Context, userId, tierId uuid.UUID) error {
	err := s.authorizer.AuthorizeAdmin(ctx)
	if err != nil {
		return err
	}
	err = s.base.AssignTier(userId, tierId)
	if err != nil {
		return err
	}
	return nil
}

func (s *authorizedService) GetUserUsage(ctx context.Context, userId uuid.UUID) (*domain.UsageSummary, error) {
	err := s.authorizer.AuthorizeAdmin(ctx)
	if err != nil {
		return nil, err
	}
	usage, err := s.base.GetUsage(userId)
	if err != nil {
		return nil, err
	}
	return usage, nil
}

func (s *authorizedService) ListInteractions(ctx context.Context, filter ports.InteractionFilter) (types.Page[domain.Interaction], error) {
	err := s.authorizer.AuthorizeAdmin(ctx)
	if err != nil {
		return types.Page[domain.Interaction]{}, err
	}
	interactions, err := s.base.ListInteractions(filter)
	if err != nil {
		return types.Page[domain.Interaction]{}, err
	}
	return interactions, nil
}
