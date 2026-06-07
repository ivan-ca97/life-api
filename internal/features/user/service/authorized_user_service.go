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
	actorId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = s.authorizer.Authorize(ctx, actorId, permissions.UsersCreate)
	if err != nil {
		return nil, err
	}
	user, err := s.base.Create(email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authorizedUserService) GetById(ctx context.Context, ownerId uuid.UUID) (*domain.User, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.UsersRead)
	if err != nil {
		return nil, err
	}
	user, err := s.base.GetById(ownerId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authorizedUserService) List(ctx context.Context, params types.PaginationParams) (types.Page[domain.User], error) {
	actorId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return types.Page[domain.User]{}, err
	}
	err = s.authorizer.Authorize(ctx, actorId, permissions.UsersRead)
	if err != nil {
		return types.Page[domain.User]{}, err
	}
	page, err := s.base.List(params)
	if err != nil {
		return types.Page[domain.User]{}, err
	}
	return page, nil
}

func (s *authorizedUserService) Update(ctx context.Context, ownerId uuid.UUID, params ports.UpdateParams) (*domain.User, error) {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.UsersUpdate)
	if err != nil {
		return nil, err
	}
	user, err := s.base.Update(ownerId, params)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authorizedUserService) Deactivate(ctx context.Context, ownerId uuid.UUID) error {
	err := s.authorizer.Authorize(ctx, ownerId, permissions.UsersDeactivate)
	if err != nil {
		return err
	}
	err = s.base.Deactivate(ownerId)
	if err != nil {
		return err
	}
	return nil
}
