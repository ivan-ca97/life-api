package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/auth/domain"
	"github.com/ivan-ca97/life/internal/features/auth/ports"
)

type sessionRepository struct {
	db *gorm.DB
}

var _ ports.SessionRepository = (*sessionRepository)(nil)

func NewSessionRepository(db *gorm.DB) *sessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) Create(session *domain.Session) error {
	err := r.db.
		Create(sessionFromDomain(session)).
		Error
	if err != nil {
		return cerr.NewInternalError("inserting session", err)
	}

	return nil
}

func (r *sessionRepository) FindById(id uuid.UUID) (*domain.Session, error) {
	var model session
	err := r.db.
		Where("id = ?", id).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSessionNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding session by id", err)
	}

	return model.toDomain(), nil
}

func (r *sessionRepository) Delete(id uuid.UUID) error {
	err := r.db.
		Where("id = ?", id).
		Delete(&session{}).
		Error
	if err != nil {
		return cerr.NewInternalError("deleting session", err)
	}

	return nil
}

func (r *sessionRepository) DeleteExpiredIfAbove(maxSessions int) error {
	var count int64
	err := r.db.
		Model(&session{}).
		Count(&count).
		Error
	if err != nil {
		return cerr.NewInternalError("counting sessions", err)
	}
	if count <= int64(maxSessions) {
		return nil
	}

	err = r.db.
		Where("expires_at < NOW()").
		Delete(&session{}).
		Error
	if err != nil {
		return cerr.NewInternalError("deleting expired sessions", err)
	}

	return nil
}
