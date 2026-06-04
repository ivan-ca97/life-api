package domain

import (
	"time"

	"github.com/google/uuid"
)

type WeightEntry struct {
	Id                uuid.UUID
	UserId            uuid.UUID
	Date              time.Time
	WeightKg          float64
	BodyFatPercentage *float64
	Notes             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
