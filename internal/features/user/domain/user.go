package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID
	Email        string
	PasswordHash string
	Active       bool
	CreatedAt    time.Time
}
