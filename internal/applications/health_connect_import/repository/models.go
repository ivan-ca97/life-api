package repository

import (
	"time"

	"github.com/google/uuid"
)

type externalHealthRecord struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId     uuid.UUID `gorm:"type:uuid;not null"`
	Source     string    `gorm:"not null"`
	Type       string    `gorm:"not null"`
	ExternalId string    `gorm:"not null"`
	RecordedAt time.Time `gorm:"not null"`
	Payload    []byte    `gorm:"type:jsonb;not null"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`
}

type syncLog struct {
	Id         int64     `gorm:"primaryKey;autoIncrement"`
	UserId     uuid.UUID `gorm:"type:uuid;not null"`
	AppVersion string    `gorm:"not null;default:''"`
	SyncedAt   time.Time `gorm:"not null"`
	ReceivedAt time.Time `gorm:"not null;autoCreateTime"`
	Result     []byte    `gorm:"type:jsonb;not null;default:'{}'"`
}
