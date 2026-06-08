package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/applications/authentication/ports"
)

type authenticationResponse struct {
	UserId    uuid.UUID `json:"user_id"`
	Token     uuid.UUID `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func authenticationResponseFromResult(result *ports.AuthenticationResult) *authenticationResponse {
	return &authenticationResponse{
		UserId:    result.UserId,
		Token:     result.Token,
		ExpiresAt: result.ExpiresAt,
	}
}
