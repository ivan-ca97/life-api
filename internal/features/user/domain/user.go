package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID
	Email        string
	Username     *string
	PasswordHash string
	GoogleId     *string
	Active       bool
	PhotoUrl     string
	HeightCm     *int
	BirthDate    *time.Time
	Sex          *string
	CreatedAt    time.Time
}
