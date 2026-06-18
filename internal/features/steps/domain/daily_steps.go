package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrStepsNotFound = errors.New("steps not found")

type DailySteps struct {
	UserId    uuid.UUID
	Date      time.Time
	Steps     int
	Source    string
	UpdatedAt time.Time
}
