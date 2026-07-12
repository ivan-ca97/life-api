package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/authorization/domain"
	"github.com/ivan-ca97/life/internal/features/authorization/ports"
)

type shareRepository struct {
	db *gorm.DB
}

var _ ports.ShareRepository = (*shareRepository)(nil)

func NewShareRepository(db *gorm.DB) *shareRepository {
	return &shareRepository{
		db: db,
	}
}

func (r *shareRepository) HasAccess(ownerId, granteeId uuid.UUID, resourceType string, needsWrite bool) (bool, error) {
	query := r.db.
		Table("shares").
		Where("owner_id = ? AND grantee_id = ? AND resource_type = ?", ownerId, granteeId, resourceType)

	if needsWrite {
		query = query.Where("can_write = true")
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, cerr.NewInternalError("checking share access", err)
	}
	return count > 0, nil
}

func (r *shareRepository) Create(s *domain.Share) error {
	model := shareFromDomain(s)
	err := r.db.Create(model).Error
	if err != nil {
		return cerr.NewInternalError("creating share", err)
	}
	return nil
}

func (r *shareRepository) ListByOwner(ownerId uuid.UUID) ([]domain.Share, error) {
	var models []share
	err := r.db.Where("owner_id = ?", ownerId).Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, cerr.NewInternalError("listing shares by owner", err)
	}
	shares := make([]domain.Share, len(models))
	for i, m := range models {
		shares[i] = *m.toDomain()
	}
	return shares, nil
}

func (r *shareRepository) ListByGrantee(granteeId uuid.UUID) ([]domain.Share, error) {
	var models []share
	err := r.db.Where("grantee_id = ?", granteeId).Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, cerr.NewInternalError("listing shares by grantee", err)
	}
	shares := make([]domain.Share, len(models))
	for i, m := range models {
		shares[i] = *m.toDomain()
	}
	return shares, nil
}

func (r *shareRepository) Update(id, ownerId uuid.UUID, canWrite bool) (*domain.Share, error) {
	result := r.db.Model(&share{}).
		Where("id = ? AND owner_id = ?", id, ownerId).
		Update("can_write", canWrite)
	if result.Error != nil {
		return nil, cerr.NewInternalError("updating share", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, cerr.NewNotFoundError("share")
	}
	var model share
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, cerr.NewInternalError("fetching updated share", err)
	}
	return model.toDomain(), nil
}

func (r *shareRepository) Delete(id, ownerId uuid.UUID) error {
	result := r.db.Where("id = ? AND owner_id = ?", id, ownerId).Delete(&share{})
	if result.Error != nil {
		return cerr.NewInternalError("deleting share", result.Error)
	}
	if result.RowsAffected == 0 {
		return cerr.NewNotFoundError("share")
	}
	return nil
}
