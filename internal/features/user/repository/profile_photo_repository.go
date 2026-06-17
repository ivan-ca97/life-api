package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
	"github.com/ivan-ca97/life/internal/features/user/ports"
)

type profilePhoto struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	Url       string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (profilePhoto) TableName() string { return "user_profile_photos" }

func (m *profilePhoto) toDomain() *domain.ProfilePhoto {
	return &domain.ProfilePhoto{
		Id:        m.Id,
		UserId:    m.UserId,
		Url:       m.Url,
		CreatedAt: m.CreatedAt,
	}
}

type profilePhotoRepository struct {
	db *gorm.DB
}

var _ ports.ProfilePhotoRepository = (*profilePhotoRepository)(nil)

func NewProfilePhotoRepository(db *gorm.DB) *profilePhotoRepository {
	return &profilePhotoRepository{db: db}
}

func (r *profilePhotoRepository) Create(photo *domain.ProfilePhoto) error {
	model := &profilePhoto{
		Id:     photo.Id,
		UserId: photo.UserId,
		Url:    photo.Url,
	}
	err := r.db.Create(model).Error
	if err != nil {
		return cerr.NewInternalError("inserting profile photo", err)
	}
	photo.CreatedAt = model.CreatedAt
	return nil
}

func (r *profilePhotoRepository) ListByUserId(userId uuid.UUID, params types.PaginationParams) (types.Page[domain.ProfilePhoto], error) {
	var models []profilePhoto
	var total int64

	q := r.db.Model(&profilePhoto{}).Where("user_id = ?", userId)

	err := q.Count(&total).Error
	if err != nil {
		return types.Page[domain.ProfilePhoto]{}, cerr.NewInternalError("counting profile photos", err)
	}

	err = q.Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&models).Error
	if err != nil {
		return types.Page[domain.ProfilePhoto]{}, cerr.NewInternalError("listing profile photos", err)
	}

	photos := make([]domain.ProfilePhoto, len(models))
	for i, m := range models {
		photos[i] = *m.toDomain()
	}

	return types.Page[domain.ProfilePhoto]{
		Items:  photos,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}, nil
}
