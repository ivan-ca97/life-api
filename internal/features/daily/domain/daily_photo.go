package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrDailyPhotoNotFound = errors.New("daily photo not found")

type DailyPhoto struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Date      time.Time
	Url       string
	Name      string
	IsPrimary bool
	CreatedAt time.Time
}
