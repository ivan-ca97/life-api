package ports

import (
	"time"

	"github.com/google/uuid"
)

type DayClosureService interface {
	Close(userId uuid.UUID, date time.Time) error
	Open(userId uuid.UUID, date time.Time) error
	IsClosed(userId uuid.UUID, date time.Time) (bool, error)
	GetClosedDates(userId uuid.UUID, from, to time.Time) (map[string]bool, error)
}
