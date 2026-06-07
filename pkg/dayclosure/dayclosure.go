package dayclosure

import (
	"time"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

type DayClosureChecker interface {
	IsClosed(userId uuid.UUID, date time.Time) (bool, error)
}

var ErrDayClosed = cerr.NewConflictError("day is closed")
