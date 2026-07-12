package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
	"github.com/ivan-ca97/life/internal/features/weight/ports"
)

type weightEntryRepository struct {
	db *gorm.DB
}

var _ ports.WeightEntryRepository = (*weightEntryRepository)(nil)

func NewWeightEntryRepository(db *gorm.DB) *weightEntryRepository {
	return &weightEntryRepository{
		db: db,
	}
}

func (r *weightEntryRepository) Create(entry *domain.WeightEntry) error {
	model := weightEntryFromDomain(entry)
	err := r.db.Create(model).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrWeightEntryConflict
		}
		return cerr.NewInternalError("inserting weight entry", err)
	}
	return nil
}

func (r *weightEntryRepository) FindById(id, userId uuid.UUID) (*domain.WeightEntry, error) {
	var model weightEntry
	err := r.db.
		Where("id = ? AND user_id = ?", id, userId).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrWeightEntryNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding weight entry by id", err)
	}
	return model.toDomain(), nil
}

func (r *weightEntryRepository) LatestByUserId(userId uuid.UUID) (*domain.WeightEntry, error) {
	var model weightEntry
	err := r.db.
		Where("user_id = ?", userId).
		Order("date DESC").
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding latest weight entry", err)
	}
	return model.toDomain(), nil
}

func (r *weightEntryRepository) List(userId uuid.UUID, params ports.ListParams) (types.Page[domain.WeightEntry], error) {
	var models []weightEntry
	var total int64

	countQuery := r.db.Model(&weightEntry{}).Where("user_id = ?", userId)
	if params.From != nil {
		countQuery = countQuery.Where("date >= ?", *params.From)
	}
	if params.To != nil {
		countQuery = countQuery.Where("date <= ?", *params.To)
	}
	err := countQuery.Count(&total).Error
	if err != nil {
		return types.Page[domain.WeightEntry]{}, cerr.NewInternalError("counting weight entries", err)
	}

	findQuery := r.db.Where("user_id = ?", userId)
	if params.From != nil {
		findQuery = findQuery.Where("date >= ?", *params.From)
	}
	if params.To != nil {
		findQuery = findQuery.Where("date <= ?", *params.To)
	}
	err = findQuery.
		Limit(params.Limit).
		Offset(params.Offset).
		Order("date DESC").
		Find(&models).
		Error
	if err != nil {
		return types.Page[domain.WeightEntry]{}, cerr.NewInternalError("listing weight entries", err)
	}

	entries := make([]domain.WeightEntry, len(models))
	for i, m := range models {
		entries[i] = *m.toDomain()
	}

	result := types.Page[domain.WeightEntry]{
		Items:  entries,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
	return result, nil
}

func (r *weightEntryRepository) Update(id, userId uuid.UUID, params ports.UpdateParams) (*domain.WeightEntry, error) {
	var count int64
	err := r.db.Model(&weightEntry{}).Where("id = ? AND user_id = ?", id, userId).Count(&count).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking weight entry existence", err)
	}
	if count == 0 {
		return nil, domain.ErrWeightEntryNotFound
	}

	updates := map[string]any{}
	if params.Date != nil {
		updates["date"] = *params.Date
	}
	if params.WeightKg != nil {
		updates["weight_kg"] = *params.WeightKg
	}
	if params.BodyFatPercentage != nil {
		updates["body_fat_percentage"] = *params.BodyFatPercentage
	}
	if params.Notes != nil {
		updates["notes"] = *params.Notes
	}

	if len(updates) > 0 {
		err = r.db.Model(&weightEntry{}).Where("id = ? AND user_id = ?", id, userId).Updates(updates).Error
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return nil, domain.ErrWeightEntryConflict
			}
			return nil, cerr.NewInternalError("updating weight entry", err)
		}
	}

	updated, err := r.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *weightEntryRepository) Delete(id, userId uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userId).Delete(&weightEntry{})
	if result.Error != nil {
		return cerr.NewInternalError("deleting weight entry", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrWeightEntryNotFound
	}
	return nil
}

func (r *weightEntryRepository) ExistsByExternalId(userId uuid.UUID, externalId string) (bool, error) {
	var count int64
	err := r.db.
		Model(&weightEntry{}).
		Where("user_id = ? AND external_id = ?", userId, externalId).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking weight entry existence by external id", err)
	}
	return count > 0, nil
}
