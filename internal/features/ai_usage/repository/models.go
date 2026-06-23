package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
)

type aiTier struct {
	Id              uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name            string    `gorm:"not null;unique"`
	MonthlyLimitUSD *float64  `gorm:"column:monthly_limit_usd"`
	IsDefault       bool      `gorm:"not null;default:false"`
	Enabled         bool      `gorm:"not null;default:true"`
	CreatedAt       time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"not null;autoUpdateTime"`
}

func (aiTier) TableName() string { return "ai_tier" }

func (m *aiTier) toDomain() domain.Tier {
	return domain.Tier{
		Id:              m.Id,
		Name:            m.Name,
		MonthlyLimitUSD: m.MonthlyLimitUSD,
		IsDefault:       m.IsDefault,
		Enabled:         m.Enabled,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

type aiUserTier struct {
	UserId       uuid.UUID `gorm:"type:uuid;primaryKey"`
	TierId       uuid.UUID `gorm:"type:uuid;not null"`
	SelfLimitUSD *float64  `gorm:"column:self_limit_usd"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
}

func (aiUserTier) TableName() string { return "ai_user_tier" }

type aiUsage struct {
	UserId        uuid.UUID `gorm:"type:uuid;primaryKey"`
	PeriodStart   time.Time `gorm:"type:date;primaryKey"`
	Requests      int       `gorm:"not null;default:0"`
	InputTokens   int64     `gorm:"not null;default:0"`
	OutputTokens  int64     `gorm:"not null;default:0"`
	CostUsdMicros int64     `gorm:"column:cost_usd_micros;not null;default:0"`
	UpdatedAt     time.Time `gorm:"not null;autoUpdateTime"`
}

func (aiUsage) TableName() string { return "ai_usage" }

func (m *aiUsage) toDomain() *domain.Usage {
	return &domain.Usage{
		UserId:       m.UserId,
		PeriodStart:  m.PeriodStart,
		Requests:     m.Requests,
		InputTokens:  m.InputTokens,
		OutputTokens: m.OutputTokens,
		CostUSD:      float64(m.CostUsdMicros) / 1_000_000,
	}
}
