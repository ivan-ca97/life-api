package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/measurements/domain"
	"github.com/ivan-ca97/life/internal/features/measurements/ports"
)

type measurementRepository struct {
	db *gorm.DB
}

var _ ports.MeasurementRepository = (*measurementRepository)(nil)

func NewMeasurementRepository(db *gorm.DB) *measurementRepository {
	return &measurementRepository{db: db}
}

func (r *measurementRepository) Upsert(m *domain.BodyMeasurement) error {
	model := &bodyMeasurement{
		Id:     m.Id,
		UserId: m.UserId,
		Date:   m.Date,
		Type:   m.Type,
		Value:  m.Value,
		Notes:  m.Notes,
	}
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}, {Name: "type"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "notes", "updated_at"}),
	}).Create(model).Error
	if err != nil {
		return cerr.NewInternalError("upserting body measurement", err)
	}
	m.CreatedAt = model.CreatedAt
	m.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *measurementRepository) FindByDate(userId uuid.UUID, date time.Time, measureType string) (*domain.BodyMeasurement, error) {
	var model bodyMeasurement
	err := r.db.
		Where("user_id = ? AND date = ? AND type = ?", userId, date, measureType).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrMeasurementNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding body measurement", err)
	}
	return model.toDomain(), nil
}

func (r *measurementRepository) List(userId uuid.UUID, params ports.ListParams) ([]domain.BodyMeasurement, error) {
	query := r.db.Where("user_id = ?", userId)
	if params.From != nil {
		query = query.Where("date >= ?", params.From)
	}
	if params.To != nil {
		query = query.Where("date <= ?", params.To)
	}
	if params.Type != nil {
		query = query.Where("type = ?", *params.Type)
	}
	var models []bodyMeasurement
	if err := query.Order("date DESC, type ASC").Find(&models).Error; err != nil {
		return nil, cerr.NewInternalError("listing body measurements", err)
	}
	results := make([]domain.BodyMeasurement, len(models))
	for i, m := range models {
		results[i] = *m.toDomain()
	}
	return results, nil
}

func (r *measurementRepository) Delete(userId uuid.UUID, date time.Time, measureType string) error {
	result := r.db.
		Where("user_id = ? AND date = ? AND type = ?", userId, date, measureType).
		Delete(&bodyMeasurement{})
	if result.Error != nil {
		return cerr.NewInternalError("deleting body measurement", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrMeasurementNotFound
	}
	return nil
}
