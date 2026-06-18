package service

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"
	"github.com/ivan-ca97/life/pkg/validate"

	"github.com/ivan-ca97/life/internal/features/user/domain"
	"github.com/ivan-ca97/life/internal/features/user/ports"
)

type userService struct {
	repository      ports.UserRepository
	photoRepository ports.ProfilePhotoRepository
}

var _ ports.UserService = (*userService)(nil)

func NewUserService(repository ports.UserRepository, photoRepository ports.ProfilePhotoRepository) *userService {
	return &userService{
		repository:      repository,
		photoRepository: photoRepository,
	}
}

func (s *userService) Create(email, password string) (*domain.User, error) {
	if err := validate.Email(email); err != nil {
		return nil, err
	}
	if err := validate.PasswordMinLength(password); err != nil {
		return nil, err
	}
	exists, err := s.repository.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, cerr.NewInternalError("hashing password", err)
	}

	user := &domain.User{
		Id:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		Active:       true,
	}
	err = s.repository.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateOAuth(email, googleId string) (*domain.User, error) {
	exists, err := s.repository.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailTaken
	}

	user := &domain.User{
		Id:       uuid.New(),
		Email:    email,
		GoogleId: &googleId,
		Active:   true,
	}
	err = s.repository.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetById(id uuid.UUID) (*domain.User, error) {
	user, err := s.repository.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByEmail(email string) (*domain.User, error) {
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) List(params types.PaginationParams) (types.Page[domain.User], error) {
	page, err := s.repository.List(params)
	if err != nil {
		return types.Page[domain.User]{}, err
	}
	return page, nil
}

func (s *userService) FindByUsername(username string) (*domain.User, error) {
	return s.repository.FindByUsername(username)
}

func (s *userService) Update(id uuid.UUID, params ports.UpdateParams) (*domain.User, error) {
	if params.Email != nil {
		if err := validate.Email(*params.Email); err != nil {
			return nil, err
		}
		exists, err := s.repository.EmailExists(*params.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrEmailTaken
		}
	}
	if params.Username != nil {
		if err := validate.NonEmpty(*params.Username, "username"); err != nil {
			return nil, err
		}
		if err := validate.MaxLength(*params.Username, 50, "username"); err != nil {
			return nil, err
		}
		exists, err := s.repository.UsernameExists(*params.Username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrUsernameTaken
		}
	}
	if params.HeightCm != nil {
		if err := validate.InRange(float64(*params.HeightCm), 50, 300, "height_cm"); err != nil {
			return nil, err
		}
	}
	if params.BirthDate != nil {
		if err := validate.NotFuture(*params.BirthDate, "birth_date"); err != nil {
			return nil, err
		}
	}
	if params.Sex != nil {
		if err := validate.OneOf(*params.Sex, []string{"male", "female", "other"}, "sex"); err != nil {
			return nil, err
		}
	}
	if params.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, cerr.NewInternalError("hashing password", err)
		}
		hashed := string(hash)
		params.Password = &hashed
	}
	user, err := s.repository.Update(id, params)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Deactivate(id uuid.UUID) error {
	_, err := s.repository.FindById(id)
	if err != nil {
		return err
	}
	err = s.repository.Deactivate(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) AddProfilePhoto(userId uuid.UUID, url string) (*domain.ProfilePhoto, error) {
	if _, err := s.repository.FindById(userId); err != nil {
		return nil, err
	}
	photo := &domain.ProfilePhoto{
		Id:     uuid.New(),
		UserId: userId,
		Url:    url,
	}
	if err := s.photoRepository.Create(photo); err != nil {
		return nil, err
	}
	if err := s.repository.UpdatePhotoUrl(userId, url); err != nil {
		return nil, err
	}
	return photo, nil
}

func (s *userService) ListProfilePhotos(userId uuid.UUID, params types.PaginationParams) (types.Page[domain.ProfilePhoto], error) {
	if _, err := s.repository.FindById(userId); err != nil {
		return types.Page[domain.ProfilePhoto]{}, err
	}
	return s.photoRepository.ListByUserId(userId, params)
}
