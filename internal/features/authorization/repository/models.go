package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/authorization/domain"
)

type role struct {
	Id   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"uniqueIndex;not null"`
}

type rolePermission struct {
	RoleId     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Permission string    `gorm:"primaryKey"`
}

type userRole struct {
	UserId uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleId uuid.UUID `gorm:"type:uuid;primaryKey"`
}

type share struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OwnerId      uuid.UUID `gorm:"type:uuid;not null"`
	GranteeId    uuid.UUID `gorm:"type:uuid;not null"`
	ResourceType string    `gorm:"not null"`
	CanWrite     bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
}

func (s *share) toDomain() *domain.Share {
	return &domain.Share{
		Id:           s.Id,
		OwnerId:      s.OwnerId,
		GranteeId:    s.GranteeId,
		ResourceType: s.ResourceType,
		CanWrite:     s.CanWrite,
		CreatedAt:    s.CreatedAt,
	}
}

func shareFromDomain(s *domain.Share) *share {
	return &share{
		Id:           s.Id,
		OwnerId:      s.OwnerId,
		GranteeId:    s.GranteeId,
		ResourceType: s.ResourceType,
		CanWrite:     s.CanWrite,
		CreatedAt:    s.CreatedAt,
	}
}
