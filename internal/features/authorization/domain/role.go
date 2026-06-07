package domain

import "github.com/google/uuid"

type Role struct {
	Id   uuid.UUID
	Name string
}
