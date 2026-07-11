package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
)

type aiModelPrice struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	Provider         string    `gorm:"not null"`
	Model            string    `gorm:"not null"`
	InputPerMillion  float64   `gorm:"column:input_per_million;not null"`
	OutputPerMillion float64   `gorm:"column:output_per_million;not null"`
	EffectiveFrom    time.Time `gorm:"not null"`
	CreatedAt        time.Time `gorm:"not null;autoCreateTime"`
}

func (aiModelPrice) TableName() string { return "ai_model_price" }

func (m *aiModelPrice) toDomain() *domain.ModelPrice {
	return &domain.ModelPrice{
		Provider:         m.Provider,
		Model:            m.Model,
		InputPerMillion:  m.InputPerMillion,
		OutputPerMillion: m.OutputPerMillion,
		EffectiveFrom:    m.EffectiveFrom,
	}
}

type aiTier struct {
	Id              uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name            string    `gorm:"not null;unique"`
	MonthlyLimitUsd *float64  `gorm:"column:monthly_limit_usd"`
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
		MonthlyLimitUsd: m.MonthlyLimitUsd,
		IsDefault:       m.IsDefault,
		Enabled:         m.Enabled,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

type aiUserTier struct {
	UserId       uuid.UUID `gorm:"type:uuid;primaryKey"`
	TierId       uuid.UUID `gorm:"type:uuid;not null"`
	SelfLimitUsd *float64  `gorm:"column:self_limit_usd"`
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
		CostUsd:      float64(m.CostUsdMicros) / 1_000_000,
	}
}

type aiInteraction struct {
	Id            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserId        uuid.UUID  `gorm:"type:uuid;not null"`
	CreatedAt     time.Time  `gorm:"not null;autoCreateTime"`
	Operation     string     `gorm:"not null"`
	Provider      string     `gorm:"not null"`
	Model         string     `gorm:"not null"`
	Status        string     `gorm:"not null"`
	ErrorType     string     `gorm:"not null;default:''"`
	InputTokens   int64      `gorm:"not null;default:0"`
	OutputTokens  int64      `gorm:"not null;default:0"`
	CostUsdMicros int64      `gorm:"column:cost_usd_micros;not null;default:0"`
	LatencyMs     int        `gorm:"not null;default:0"`
	ProviderCalls int        `gorm:"not null;default:0"`
	CorrelationId *uuid.UUID `gorm:"type:uuid"`
	InputSummary  string     `gorm:"not null;default:''"`
	Metadata      []byte     `gorm:"type:jsonb;not null;default:'{}'"`
}

func (aiInteraction) TableName() string { return "ai_interaction" }

func (m *aiInteraction) toDomain() domain.Interaction {
	return domain.Interaction{
		Id:            m.Id,
		UserId:        m.UserId,
		CreatedAt:     m.CreatedAt,
		Operation:     m.Operation,
		Provider:      m.Provider,
		Model:         m.Model,
		Status:        m.Status,
		ErrorType:     m.ErrorType,
		InputTokens:   m.InputTokens,
		OutputTokens:  m.OutputTokens,
		CostUsd:       float64(m.CostUsdMicros) / 1_000_000,
		LatencyMs:     m.LatencyMs,
		ProviderCalls: m.ProviderCalls,
		CorrelationId: m.CorrelationId,
		InputSummary:  m.InputSummary,
		Metadata:      json.RawMessage(m.Metadata),
	}
}
