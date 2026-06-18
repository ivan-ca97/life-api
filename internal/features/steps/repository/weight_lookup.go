package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/internal/features/steps/ports"
)

type weightLookup struct {
	db *gorm.DB
}

var _ ports.WeightLookup = (*weightLookup)(nil)

func NewWeightLookup(db *gorm.DB) *weightLookup {
	return &weightLookup{db: db}
}

func (w *weightLookup) LatestWeightKg(userId uuid.UUID) (float64, bool, error) {
	var kg float64
	err := w.db.
		Table("weight_entries").
		Select("weight_kg").
		Where("user_id = ?", userId).
		Order("date DESC").
		Limit(1).
		Scan(&kg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || kg == 0 {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return kg, true, nil
}
