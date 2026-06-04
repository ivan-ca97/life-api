package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
)

type goal struct {
	Id                   uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId               uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	DailyCalories        *float64
	DailyProteinGrams    *float64
	DailyCarbsGrams      *float64
	DailyFatGrams        *float64
	DailyFiberGrams      *float64
	DailySteps           *int
	DailyExerciseMinutes *int
	TargetWeightKg       *float64
	StartedAt            time.Time `gorm:"not null"`
	CreatedAt            time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"not null;autoUpdateTime"`
}

func (m *goal) toDomain() *domain.Goal {
	return &domain.Goal{
		Id:                   m.Id,
		UserId:               m.UserId,
		DailyCalories:        m.DailyCalories,
		DailyProteinGrams:    m.DailyProteinGrams,
		DailyCarbsGrams:      m.DailyCarbsGrams,
		DailyFatGrams:        m.DailyFatGrams,
		DailyFiberGrams:      m.DailyFiberGrams,
		DailySteps:           m.DailySteps,
		DailyExerciseMinutes: m.DailyExerciseMinutes,
		TargetWeightKg:       m.TargetWeightKg,
		StartedAt:            m.StartedAt,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func goalFromDomain(g *domain.Goal) *goal {
	return &goal{
		Id:                   g.Id,
		UserId:               g.UserId,
		DailyCalories:        g.DailyCalories,
		DailyProteinGrams:    g.DailyProteinGrams,
		DailyCarbsGrams:      g.DailyCarbsGrams,
		DailyFatGrams:        g.DailyFatGrams,
		DailyFiberGrams:      g.DailyFiberGrams,
		DailySteps:           g.DailySteps,
		DailyExerciseMinutes: g.DailyExerciseMinutes,
		TargetWeightKg:       g.TargetWeightKg,
		StartedAt:            g.StartedAt,
		CreatedAt:            g.CreatedAt,
		UpdatedAt:            g.UpdatedAt,
	}
}
