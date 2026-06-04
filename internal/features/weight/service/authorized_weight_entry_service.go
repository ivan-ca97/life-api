package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
	"github.com/ivan-ca97/life/internal/features/weight/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedWeightEntryService struct {
	base       ports.WeightEntryService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedWeightEntryService = (*authorizedWeightEntryService)(nil)

func NewAuthorizedWeightEntryService(base ports.WeightEntryService, authorizer auth.AuthorizationService) *authorizedWeightEntryService {
	return &authorizedWeightEntryService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedWeightEntryService) Create(ctx context.Context, params ports.CreateParams) (*domain.WeightEntry, error) {
	err := s.authorizer.Require(ctx, permissions.WeightCreate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	entry, err := s.base.Create(userId, params)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *authorizedWeightEntryService) GetById(ctx context.Context, id uuid.UUID) (*domain.WeightEntry, error) {
	err := s.authorizer.Require(ctx, permissions.WeightRead)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	entry, err := s.base.GetById(id, userId)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *authorizedWeightEntryService) List(ctx context.Context, params ports.ListParams) (types.Page[domain.WeightEntry], error) {
	err := s.authorizer.Require(ctx, permissions.WeightRead)
	if err != nil {
		return types.Page[domain.WeightEntry]{}, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return types.Page[domain.WeightEntry]{}, err
	}
	page, err := s.base.List(userId, params)
	if err != nil {
		return types.Page[domain.WeightEntry]{}, err
	}
	return page, nil
}

func (s *authorizedWeightEntryService) Update(ctx context.Context, id uuid.UUID, params ports.UpdateParams) (*domain.WeightEntry, error) {
	err := s.authorizer.Require(ctx, permissions.WeightUpdate)
	if err != nil {
		return nil, err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	entry, err := s.base.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *authorizedWeightEntryService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.authorizer.Require(ctx, permissions.WeightDelete)
	if err != nil {
		return err
	}
	userId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return err
	}
	err = s.base.Delete(id, userId)
	if err != nil {
		return err
	}
	return nil
}
