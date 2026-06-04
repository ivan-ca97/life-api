package auth

import "context"

type binaryAuthorizationService struct{}

var _ AuthorizationService = (*binaryAuthorizationService)(nil)

func NewBinaryAuthorizationService() *binaryAuthorizationService {
	return &binaryAuthorizationService{}
}

func (s *binaryAuthorizationService) Require(ctx context.Context, _ string) error {
	_, err := ActorFromContext(ctx)
	if err != nil {
		return ErrForbidden
	}
	return nil
}

func (s *binaryAuthorizationService) RequireOn(ctx context.Context, _ string, _ any) error {
	_, err := ActorFromContext(ctx)
	if err != nil {
		return ErrForbidden
	}
	return nil
}
