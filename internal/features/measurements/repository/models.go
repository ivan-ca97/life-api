package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/measurements/domain"
)

type bodyMeasurement struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	Date      time.Time `gorm:"type:date;not null"`
	Type      string    `gorm:"not null"`
	Value     float64   `gorm:"not null"`
	Notes     string    `gorm:"not null;default:''"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (bodyMeasurement) TableName() string { return "body_measurements" }

func (m *bodyMeasurement) toDomain() *domain.BodyMeasurement {
	return &domain.BodyMeasurement{
		Id:        m.Id,
		UserId:    m.UserId,
		Date:      m.Date,
		Type:      m.Type,
		Value:     m.Value,
		Notes:     m.Notes,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
