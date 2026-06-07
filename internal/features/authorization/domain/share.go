package domain

import (
	"time"

	"github.com/google/uuid"
)

var ValidResourceTypes = map[string]bool{
	"meals":     true,
	"exercises": true,
	"weight":    true,
	"foods":     true,
	"goals":     true,
	"daily":     true,
}

type Share struct {
	Id           uuid.UUID
	OwnerId      uuid.UUID
	GranteeId    uuid.UUID
	ResourceType string
	CanWrite     bool
	CreatedAt    time.Time
}
