package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
	"github.com/ivan-ca97/life/internal/features/goal/ports"
)

type goalRepository struct {
	db *gorm.DB
}

var _ ports.GoalRepository = (*goalRepository)(nil)

func NewGoalRepository(db *gorm.DB) *goalRepository {
	return &goalRepository{
		db: db,
	}
}

func (r *goalRepository) FindByUserId(userId uuid.UUID) (*domain.Goal, error) {
	var model goal
	err := r.db.
		Where("user_id = ?", userId).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrGoalNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding goal by user id", err)
	}
	return model.toDomain(), nil
}

func (r *goalRepository) Upsert(g *domain.Goal) (*domain.Goal, error) {
	model := goalFromDomain(g)
	err := r.db.Save(model).Error
	if err != nil {
		return nil, cerr.NewInternalError("upserting goal", err)
	}
	return model.toDomain(), nil
}
