package domain

import (
	"time"

	"github.com/google/uuid"
)

type DayClosure struct {
	UserId   uuid.UUID
	Date     time.Time
	ClosedAt time.Time
}
