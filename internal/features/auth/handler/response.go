package handler

import (
	"time"

	"github.com/google/uuid"
)

type loginResponse struct {
	Token     uuid.UUID `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
