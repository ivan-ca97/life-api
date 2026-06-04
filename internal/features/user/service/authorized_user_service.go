package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
	"github.com/ivan-ca97/life/internal/features/user/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedUserService struct {
	base       ports.UserService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedUserService = (*authorizedUserService)(nil)

func NewAuthorizedUserService(base ports.UserService, authorizer auth.AuthorizationService) *authorizedUserService {
	return &authorizedUserService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedUserService) Create(ctx context.Context, email, password string) (*domain.User, error) {
	err := s.authorizer.Require(ctx, permissions.UsersCreate)
	if err != nil {
		return nil, err
	}
	user, err := s.base.Create(email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authorizedUserService) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	err := s.authorizer.Require(ctx, permissions.UsersRead)
	if err != nil {
		return nil, err
	}
	user, err := s.base.GetById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authorizedUserService) List(ctx context.Context, params types.PaginationParams) (types.Page[domain.User], error) {
	err := s.authorizer.Require(ctx, permissions.UsersRead)
	if err != nil {
		return types.Page[domain.User]{}, err
	}
	page, err := s.base.List(params)
	if err != nil {
		return types.Page[domain.User]{}, err
	}
	return page, nil
}

func (s *authorizedUserService) Update(ctx context.Context, id uuid.UUID, params ports.UpdateParams) (*domain.User, error) {
	err := s.authorizer.Require(ctx, permissions.UsersUpdate)
	if err != nil {
		return nil, err
	}
	user, err := s.base.Update(id, params)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authorizedUserService) Deactivate(ctx context.Context, id uuid.UUID) error {
	err := s.authorizer.Require(ctx, permissions.UsersDeactivate)
	if err != nil {
		return err
	}
	err = s.base.Deactivate(id)
	if err != nil {
		return err
	}
	return nil
}
