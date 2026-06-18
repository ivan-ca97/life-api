package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/steps/domain"
	"github.com/ivan-ca97/life/internal/features/steps/ports"
)

type stepsRepository struct {
	db *gorm.DB
}

var _ ports.StepsRepository = (*stepsRepository)(nil)

func NewStepsRepository(db *gorm.DB) *stepsRepository {
	return &stepsRepository{db: db}
}

func (r *stepsRepository) Upsert(entry *domain.DailySteps) error {
	model := &dailySteps{
		UserId: entry.UserId,
		Date:   entry.Date,
		Steps:  entry.Steps,
		Source: entry.Source,
	}
	err := r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{"steps", "source", "updated_at"}),
		}).
		Create(model).Error
	if err != nil {
		return cerr.NewInternalError("upserting daily steps", err)
	}
	entry.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *stepsRepository) FindByDate(userId uuid.UUID, date time.Time) (*domain.DailySteps, error) {
	var model dailySteps
	err := r.db.
		Where("user_id = ? AND date = ?", userId, date).
		First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrStepsNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding daily steps", err)
	}
	return model.toDomain(), nil
}

func (r *stepsRepository) List(userId uuid.UUID, params ports.ListParams) ([]domain.DailySteps, error) {
	var models []dailySteps
	q := r.db.Where("user_id = ?", userId)
	if params.From != nil {
		q = q.Where("date >= ?", *params.From)
	}
	if params.To != nil {
		q = q.Where("date <= ?", *params.To)
	}
	err := q.Order("date DESC").Find(&models).Error
	if err != nil {
		return nil, cerr.NewInternalError("listing daily steps", err)
	}
	result := make([]domain.DailySteps, len(models))
	for i, m := range models {
		result[i] = *m.toDomain()
	}
	return result, nil
}

func (r *stepsRepository) Delete(userId uuid.UUID, date time.Time) error {
	res := r.db.
		Where("user_id = ? AND date = ?", userId, date).
		Delete(&dailySteps{})
	if res.Error != nil {
		return cerr.NewInternalError("deleting daily steps", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrStepsNotFound
	}
	return nil
}
