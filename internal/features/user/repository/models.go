package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/user/domain"
)

type user struct {
	Id           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Email        string     `gorm:"uniqueIndex;not null"`
	PasswordHash string     `gorm:"not null"`
	Active       bool       `gorm:"not null;default:true"`
	HeightCm     *int       `gorm:""`
	BirthDate    *time.Time `gorm:"type:date"`
	Sex          *string    `gorm:""`
	CreatedAt    time.Time  `gorm:"not null;autoCreateTime"`
}

func (m *user) toDomain() *domain.User {
	return &domain.User{
		Id:           m.Id,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Active:       m.Active,
		HeightCm:     m.HeightCm,
		BirthDate:    m.BirthDate,
		Sex:          m.Sex,
		CreatedAt:    m.CreatedAt,
	}
}

func userFromDomain(u *domain.User) *user {
	return &user{
		Id:           u.Id,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Active:       u.Active,
		HeightCm:     u.HeightCm,
		BirthDate:    u.BirthDate,
		Sex:          u.Sex,
		CreatedAt:    u.CreatedAt,
	}
}
