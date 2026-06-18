package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/steps/domain"
)

type dailySteps struct {
	UserId    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Date      time.Time `gorm:"type:date;primaryKey"`
	Steps     int       `gorm:"not null"`
	Source    string    `gorm:"not null;default:''"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (dailySteps) TableName() string { return "daily_steps" }

func (m *dailySteps) toDomain() *domain.DailySteps {
	return &domain.DailySteps{
		UserId:    m.UserId,
		Date:      m.Date,
		Steps:     m.Steps,
		Source:    m.Source,
		UpdatedAt: m.UpdatedAt,
	}
}
