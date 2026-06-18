package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

var ErrMeasurementNotFound = cerr.NewNotFoundError("body measurement")

var ErrDuplicateMeasurement = errors.New("measurement already exists for this date and type")

type BodyMeasurement struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Date      time.Time
	Type      string
	Value     float64
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
