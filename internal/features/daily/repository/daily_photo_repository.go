package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type dailyPhotoModel struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	Date      time.Time `gorm:"type:date;not null"`
	Url       string    `gorm:"not null"`
	Name      string    `gorm:"not null;default:''"`
	IsPrimary bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (dailyPhotoModel) TableName() string {
	return "daily_photos"
}

func (m *dailyPhotoModel) toDomain() *domain.DailyPhoto {
	return &domain.DailyPhoto{
		Id:        m.Id,
		UserId:    m.UserId,
		Date:      m.Date,
		Url:       m.Url,
		Name:      m.Name,
		IsPrimary: m.IsPrimary,
		CreatedAt: m.CreatedAt,
	}
}

type dailyPhotoRepository struct {
	db *gorm.DB
}

var _ ports.PhotoRepository = (*dailyPhotoRepository)(nil)

func NewDailyPhotoRepository(db *gorm.DB) *dailyPhotoRepository {
	return &dailyPhotoRepository{db: db}
}

func (r *dailyPhotoRepository) Create(photo *domain.DailyPhoto) error {
	model := &dailyPhotoModel{
		Id:        photo.Id,
		UserId:    photo.UserId,
		Date:      photo.Date,
		Url:       photo.Url,
		Name:      photo.Name,
		IsPrimary: photo.IsPrimary,
	}
	err := r.db.Create(model).Error
	if err != nil {
		return cerr.NewInternalError("inserting daily photo", err)
	}
	photo.CreatedAt = model.CreatedAt
	return nil
}

func (r *dailyPhotoRepository) FindById(id, userId uuid.UUID) (*domain.DailyPhoto, error) {
	var model dailyPhotoModel
	err := r.db.Where("id = ? AND user_id = ?", id, userId).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrDailyPhotoNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding daily photo", err)
	}
	return model.toDomain(), nil
}

func (r *dailyPhotoRepository) ListByDate(userId uuid.UUID, date time.Time) ([]domain.DailyPhoto, error) {
	var models []dailyPhotoModel
	err := r.db.
		Where("user_id = ? AND date = ?", userId, date).
		Order("is_primary DESC, created_at ASC").
		Find(&models).Error
	if err != nil {
		return nil, cerr.NewInternalError("listing daily photos", err)
	}
	photos := make([]domain.DailyPhoto, len(models))
	for i, m := range models {
		photos[i] = *m.toDomain()
	}
	return photos, nil
}

func (r *dailyPhotoRepository) Update(id, userId uuid.UUID, name *string, isPrimary *bool) (*domain.DailyPhoto, error) {
	updates := map[string]any{}
	if name != nil {
		updates["name"] = *name
	}
	if isPrimary != nil {
		updates["is_primary"] = *isPrimary
	}
	if len(updates) > 0 {
		err := r.db.Model(&dailyPhotoModel{}).Where("id = ? AND user_id = ?", id, userId).Updates(updates).Error
		if err != nil {
			return nil, cerr.NewInternalError("updating daily photo", err)
		}
	}
	return r.FindById(id, userId)
}

func (r *dailyPhotoRepository) UnsetPrimary(userId uuid.UUID, date time.Time) error {
	err := r.db.Model(&dailyPhotoModel{}).
		Where("user_id = ? AND date = ? AND is_primary = true", userId, date).
		Update("is_primary", false).Error
	if err != nil {
		return cerr.NewInternalError("unsetting primary daily photo", err)
	}
	return nil
}

func (r *dailyPhotoRepository) Delete(id, userId uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userId).Delete(&dailyPhotoModel{})
	if result.Error != nil {
		return cerr.NewInternalError("deleting daily photo", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrDailyPhotoNotFound
	}
	return nil
}
