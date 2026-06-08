package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/authentication/domain"
	"github.com/ivan-ca97/life/internal/features/authentication/ports"
)

type authenticationService struct {
	sessionRepository ports.SessionRepository
}

var _ ports.AuthenticationService = (*authenticationService)(nil)

func NewAuthenticationService(sessionRepository ports.SessionRepository) *authenticationService {
	return &authenticationService{
		sessionRepository: sessionRepository,
	}
}

func (s *authenticationService) CreateSession(userId uuid.UUID) (*domain.Session, error) {
	session := &domain.Session{
		Id:        uuid.New(),
		UserId:    userId,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err := s.sessionRepository.Create(session)
	if err != nil {
		return nil, err
	}

	s.sessionRepository.DeleteExpiredIfAbove(100)
	return session, nil
}

func (s *authenticationService) Validate(sessionId uuid.UUID) (*domain.Session, error) {
	session, err := s.sessionRepository.FindById(sessionId)
	if err != nil {
		return nil, err
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, domain.ErrSessionExpired
	}

	return session, nil
}

func (s *authenticationService) Logout(sessionId uuid.UUID) error {
	err := s.sessionRepository.Delete(sessionId)
	if err != nil {
		return err
	}

	return nil
}
