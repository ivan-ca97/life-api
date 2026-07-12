package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

const (
	correctionMealType     = "correction"
	correctionMealName     = "Daily Correction"
	correctionExerciseType = "manual_adjustment"
	correctionExerciseName = "Daily Correction"
)

type correctionMeal struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId       uuid.UUID `gorm:"type:uuid;not null"`
	Date         time.Time `gorm:"type:date;not null"`
	Type         string    `gorm:"not null"`
	Name         string    `gorm:"default:''"`
	PhotoUrl     string    `gorm:"not null;default:''"`
	Calories     *float64
	ProteinGrams *float64
	CarbsGrams   *float64
	FatGrams     *float64
	FiberGrams   *float64
	Notes        string    `gorm:"not null;default:''"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
}

func (correctionMeal) TableName() string { return "meals" }

type correctionExercise struct {
	Id                      uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId                  uuid.UUID `gorm:"type:uuid;not null"`
	Date                    time.Time `gorm:"type:date;not null"`
	Type                    string    `gorm:"not null"`
	Name                    string    `gorm:"not null"`
	EstimatedCaloriesBurned *float64
	Steps                   *int
	DurationSeconds         *int
	DistanceMeters          *float64
	Notes                   string    `gorm:"not null;default:''"`
	CreatedAt               time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt               time.Time `gorm:"not null;autoUpdateTime"`
}

func (correctionExercise) TableName() string { return "exercises" }

type correctionRepository struct {
	db *gorm.DB
}

var _ ports.CorrectionRepository = (*correctionRepository)(nil)

func NewCorrectionRepository(db *gorm.DB) *correctionRepository {
	return &correctionRepository{
		db: db,
	}
}

func (r *correctionRepository) GetCorrection(userId uuid.UUID, date time.Time) (*domain.Correction, error) {
	var meal correctionMeal
	mealErr := r.db.
		Where("user_id = ? AND date = ? AND type = ?", userId, date, correctionMealType).
		Take(&meal).Error
	hasMeal := true
	if errors.Is(mealErr, gorm.ErrRecordNotFound) {
		hasMeal = false
	} else if mealErr != nil {
		return nil, cerr.NewInternalError("fetching correction meal", mealErr)
	}

	var ex correctionExercise
	exErr := r.db.
		Where("user_id = ? AND date = ? AND type = ? AND name = ?", userId, date, correctionExerciseType, correctionExerciseName).
		Take(&ex).Error
	hasExercise := true
	if errors.Is(exErr, gorm.ErrRecordNotFound) {
		hasExercise = false
	} else if exErr != nil {
		return nil, cerr.NewInternalError("fetching correction exercise", exErr)
	}

	if !hasMeal && !hasExercise {
		return nil, nil
	}

	correction := &domain.Correction{
		Date: date,
	}
	if hasMeal {
		correction.Calories = meal.Calories
		correction.ProteinGrams = meal.ProteinGrams
		correction.CarbsGrams = meal.CarbsGrams
		correction.FatGrams = meal.FatGrams
		correction.FiberGrams = meal.FiberGrams
		correction.Notes = meal.Notes
	}
	if hasExercise {
		correction.CaloriesBurned = ex.EstimatedCaloriesBurned
		correction.Steps = ex.Steps
		correction.DurationSeconds = ex.DurationSeconds
		correction.DistanceMeters = ex.DistanceMeters
		if !hasMeal {
			correction.Notes = ex.Notes
		}
	}
	return correction, nil
}

func (r *correctionRepository) UpsertCorrection(userId uuid.UUID, correction *domain.Correction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := r.upsertMealCorrection(tx, userId, correction)
		if err != nil {
			return err
		}
		err = r.upsertExerciseCorrection(tx, userId, correction)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *correctionRepository) DeleteCorrection(userId uuid.UUID, date time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where("user_id = ? AND date = ? AND type = ?", userId, date, correctionMealType).
			Delete(&correctionMeal{}).Error
		if err != nil {
			return cerr.NewInternalError("deleting correction meal", err)
		}
		err = tx.
			Where("user_id = ? AND date = ? AND type = ? AND name = ?", userId, date, correctionExerciseType, correctionExerciseName).
			Delete(&correctionExercise{}).Error
		if err != nil {
			return cerr.NewInternalError("deleting correction exercise", err)
		}
		return nil
	})
}

func (r *correctionRepository) upsertMealCorrection(tx *gorm.DB, userId uuid.UUID, correction *domain.Correction) error {
	if !correction.HasMealFields() {
		return tx.
			Where("user_id = ? AND date = ? AND type = ?", userId, correction.Date, correctionMealType).
			Delete(&correctionMeal{}).Error
	}

	var existing correctionMeal
	err := tx.
		Where("user_id = ? AND date = ? AND type = ?", userId, correction.Date, correctionMealType).
		Take(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		m := correctionMeal{
			Id:           uuid.New(),
			UserId:       userId,
			Date:         correction.Date,
			Type:         correctionMealType,
			Name:         correctionMealName,
			Calories:     correction.Calories,
			ProteinGrams: correction.ProteinGrams,
			CarbsGrams:   correction.CarbsGrams,
			FatGrams:     correction.FatGrams,
			FiberGrams:   correction.FiberGrams,
			Notes:        correction.Notes,
		}
		err := tx.Create(&m).Error
		if err != nil {
			return cerr.NewInternalError("creating correction meal", err)
		}
		return nil
	}
	if err != nil {
		return cerr.NewInternalError("fetching correction meal for upsert", err)
	}

	existing.Calories = correction.Calories
	existing.ProteinGrams = correction.ProteinGrams
	existing.CarbsGrams = correction.CarbsGrams
	existing.FatGrams = correction.FatGrams
	existing.FiberGrams = correction.FiberGrams
	existing.Notes = correction.Notes
	err = tx.Save(&existing).Error
	if err != nil {
		return cerr.NewInternalError("updating correction meal", err)
	}
	return nil
}

func (r *correctionRepository) upsertExerciseCorrection(tx *gorm.DB, userId uuid.UUID, correction *domain.Correction) error {
	if !correction.HasExerciseFields() {
		return tx.
			Where("user_id = ? AND date = ? AND type = ? AND name = ?", userId, correction.Date, correctionExerciseType, correctionExerciseName).
			Delete(&correctionExercise{}).Error
	}

	var existing correctionExercise
	err := tx.
		Where("user_id = ? AND date = ? AND type = ? AND name = ?", userId, correction.Date, correctionExerciseType, correctionExerciseName).
		Take(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		e := correctionExercise{
			Id:                      uuid.New(),
			UserId:                  userId,
			Date:                    correction.Date,
			Type:                    correctionExerciseType,
			Name:                    correctionExerciseName,
			EstimatedCaloriesBurned: correction.CaloriesBurned,
			Steps:                   correction.Steps,
			DurationSeconds:         correction.DurationSeconds,
			DistanceMeters:          correction.DistanceMeters,
			Notes:                   correction.Notes,
		}
		err := tx.Create(&e).Error
		if err != nil {
			return cerr.NewInternalError("creating correction exercise", err)
		}
		return nil
	}
	if err != nil {
		return cerr.NewInternalError("fetching correction exercise for upsert", err)
	}

	existing.EstimatedCaloriesBurned = correction.CaloriesBurned
	existing.Steps = correction.Steps
	existing.DurationSeconds = correction.DurationSeconds
	existing.DistanceMeters = correction.DistanceMeters
	existing.Notes = correction.Notes
	err = tx.Save(&existing).Error
	if err != nil {
		return cerr.NewInternalError("updating correction exercise", err)
	}
	return nil
}
