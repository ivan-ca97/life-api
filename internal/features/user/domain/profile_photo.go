package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProfilePhoto struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Url       string
	CreatedAt time.Time
}
