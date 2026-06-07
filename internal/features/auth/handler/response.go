package handler

import (
	"time"

	"github.com/google/uuid"
)

type loginResponse struct {
	UserId    uuid.UUID `json:"user_id"`
	Token     uuid.UUID `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type registerResponse struct {
	UserId    uuid.UUID `json:"user_id"`
	Token     uuid.UUID `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
