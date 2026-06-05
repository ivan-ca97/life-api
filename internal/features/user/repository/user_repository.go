package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
	"github.com/ivan-ca97/life/internal/features/user/ports"
)

type userRepository struct {
	db *gorm.DB
}

var _ ports.UserRepository = (*userRepository)(nil)

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	err := r.db.
		Create(userFromDomain(user)).
		Error
	if err != nil {
		return cerr.NewInternalError("inserting user", err)
	}
	return nil
}

func (r *userRepository) FindById(id uuid.UUID) (*domain.User, error) {
	var model user
	err := r.db.
		Where("id = ?", id).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding user by id", err)
	}
	return model.toDomain(), nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var model user
	err := r.db.
		Where("email = ?", email).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding user by email", err)
	}
	return model.toDomain(), nil
}

func (r *userRepository) List(params types.PaginationParams) (types.Page[domain.User], error) {
	var models []user
	var total int64

	err := r.db.
		Model(&user{}).
		Count(&total).
		Error
	if err != nil {
		return types.Page[domain.User]{}, cerr.NewInternalError("counting users", err)
	}

	err = r.db.
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&models).
		Error
	if err != nil {
		return types.Page[domain.User]{}, cerr.NewInternalError("listing users", err)
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = *m.toDomain()
	}

	result := types.Page[domain.User]{
		Items:  users,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
	return result, nil
}

func (r *userRepository) Update(id uuid.UUID, params ports.UpdateParams) (*domain.User, error) {
	updates := map[string]any{}
	if params.Email != nil {
		updates["email"] = *params.Email
	}
	if params.Password != nil {
		updates["password_hash"] = *params.Password
	}
	if params.HeightCm != nil {
		updates["height_cm"] = *params.HeightCm
	}
	if params.BirthDate != nil {
		updates["birth_date"] = *params.BirthDate
	}
	if params.Sex != nil {
		updates["sex"] = *params.Sex
	}

	err := r.db.
		Model(&user{}).
		Where("id = ?", id).
		Updates(updates).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("updating user", err)
	}

	user, err := r.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Deactivate(id uuid.UUID) error {
	err := r.db.
		Model(&user{}).
		Where("id = ?", id).
		Update("active", false).
		Error
	if err != nil {
		return cerr.NewInternalError("deactivating user", err)
	}
	return nil
}

func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.
		Model(&user{}).
		Where("email = ?", email).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking email existence", err)
	}
	return count > 0, nil
}
