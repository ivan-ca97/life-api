package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
)

type weightEntry struct {
	Id                uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId            uuid.UUID `gorm:"type:uuid;not null"`
	Date              time.Time `gorm:"type:date;not null"`
	WeightKg          float64   `gorm:"not null"`
	BodyFatPercentage *float64
	Notes             string `gorm:"not null;default:''"`
	ExternalId        *string
	CreatedAt         time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"not null;autoUpdateTime"`
}

func (m *weightEntry) toDomain() *domain.WeightEntry {
	return &domain.WeightEntry{
		Id:                m.Id,
		UserId:            m.UserId,
		Date:              m.Date,
		WeightKg:          m.WeightKg,
		BodyFatPercentage: m.BodyFatPercentage,
		Notes:             m.Notes,
		ExternalId:        m.ExternalId,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func weightEntryFromDomain(e *domain.WeightEntry) *weightEntry {
	return &weightEntry{
		Id:                e.Id,
		UserId:            e.UserId,
		Date:              e.Date,
		WeightKg:          e.WeightKg,
		BodyFatPercentage: e.BodyFatPercentage,
		Notes:             e.Notes,
		ExternalId:        e.ExternalId,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}
