package service

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ivan-ca97/life/internal/features/auth/domain"
	"github.com/ivan-ca97/life/internal/features/auth/ports"
	user_domain "github.com/ivan-ca97/life/internal/features/user/domain"
	user_ports "github.com/ivan-ca97/life/internal/features/user/ports"
)

type authService struct {
	sessionRepository ports.SessionRepository
	userService       user_ports.UserService
}

var _ ports.AuthService = (*authService)(nil)

func NewAuthService(sessionRepository ports.SessionRepository, userService user_ports.UserService) *authService {
	return &authService{
		sessionRepository: sessionRepository,
		userService:       userService,
	}
}

func (s *authService) Login(email, password string) (*domain.Session, error) {
	user, err := s.userService.GetByEmail(email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.Active {
		return nil, user_domain.ErrUserInactive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	session := &domain.Session{
		Id:        uuid.New(),
		UserId:    user.Id,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err = s.sessionRepository.Create(session)
	if err != nil {
		return nil, err
	}

	s.sessionRepository.DeleteExpiredIfAbove(100)
	return session, nil
}

func (s *authService) CreateSession(userId uuid.UUID) (*domain.Session, error) {
	session := &domain.Session{
		Id:        uuid.New(),
		UserId:    userId,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err := s.sessionRepository.Create(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *authService) Validate(sessionId uuid.UUID) (*domain.Session, error) {
	session, err := s.sessionRepository.FindById(sessionId)
	if err != nil {
		return nil, err
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, domain.ErrSessionExpired
	}

	return session, nil
}

func (s *authService) Logout(sessionId uuid.UUID) error {
	err := s.sessionRepository.Delete(sessionId)
	if err != nil {
		return err
	}

	return nil
}
