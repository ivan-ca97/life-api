package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type dailyPhotoService struct {
	repository ports.PhotoRepository
}

var _ ports.PhotoService = (*dailyPhotoService)(nil)

func NewDailyPhotoService(repository ports.PhotoRepository) *dailyPhotoService {
	return &dailyPhotoService{
		repository: repository,
	}
}

func (s *dailyPhotoService) Create(userId uuid.UUID, params ports.CreatePhotoParams) (*domain.DailyPhoto, error) {
	if params.IsPrimary {
		err := s.repository.UnsetPrimary(userId, params.Date)
		if err != nil {
			return nil, err
		}
	}
	photo := &domain.DailyPhoto{
		Id:        uuid.New(),
		UserId:    userId,
		Date:      params.Date,
		Url:       params.Url,
		Name:      params.Name,
		IsPrimary: params.IsPrimary,
	}
	err := s.repository.Create(photo)
	if err != nil {
		return nil, err
	}
	return photo, nil
}

func (s *dailyPhotoService) List(userId uuid.UUID, date time.Time) ([]domain.DailyPhoto, error) {
	photos, err := s.repository.ListByDate(userId, date)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func (s *dailyPhotoService) Update(id, userId uuid.UUID, params ports.UpdatePhotoParams) (*domain.DailyPhoto, error) {
	if params.IsPrimary != nil && *params.IsPrimary {
		photo, err := s.repository.FindById(id, userId)
		if err != nil {
			return nil, err
		}
		err = s.repository.UnsetPrimary(userId, photo.Date)
		if err != nil {
			return nil, err
		}
	}
	photo, err := s.repository.Update(id, userId, params.Name, params.IsPrimary)
	if err != nil {
		return nil, err
	}
	return photo, nil
}

func (s *dailyPhotoService) Delete(id, userId uuid.UUID) error {
	err := s.repository.Delete(id, userId)
	if err != nil {
		return err
	}
	return nil
}
