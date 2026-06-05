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
	HeightCm     *int
	BirthDate    *time.Time
	Sex          *string
	CreatedAt    time.Time
}
