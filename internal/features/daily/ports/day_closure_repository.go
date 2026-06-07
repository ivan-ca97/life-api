package ports

import (
	"time"

	"github.com/google/uuid"
)

type DayClosureRepository interface {
	IsClosed(userId uuid.UUID, date time.Time) (bool, error)
	Close(userId uuid.UUID, date time.Time) error
	Open(userId uuid.UUID, date time.Time) error
	GetClosedDates(userId uuid.UUID, from, to time.Time) (map[string]bool, error)
}
