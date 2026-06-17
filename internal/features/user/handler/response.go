package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
)

type profilePhotoResponse struct {
	Id        uuid.UUID `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

func profilePhotoFromDomain(p *domain.ProfilePhoto) *profilePhotoResponse {
	return &profilePhotoResponse{
		Id:        p.Id,
		Url:       p.Url,
		CreatedAt: p.CreatedAt,
	}
}

type profilePhotoPage struct {
	Items  []profilePhotoResponse `json:"items"`
	Total  int64                  `json:"total"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}

func newProfilePhotoPage(page types.Page[domain.ProfilePhoto]) *profilePhotoPage {
	items := make([]profilePhotoResponse, len(page.Items))
	for i, p := range page.Items {
		items[i] = *profilePhotoFromDomain(&p)
	}
	return &profilePhotoPage{
		Items:  items,
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
}

type userResponse struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	PhotoUrl  string    `json:"photo_url"`
	HeightCm  *int      `json:"height_cm,omitempty"`
	BirthDate *string   `json:"birth_date,omitempty"`
	Sex       *string   `json:"sex,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func userFromDomain(u *domain.User) *userResponse {
	var birthDate *string
	if u.BirthDate != nil {
		s := u.BirthDate.Format("2006-01-02")
		birthDate = &s
	}
	return &userResponse{
		Id:        u.Id,
		Email:     u.Email,
		Active:    u.Active,
		PhotoUrl:  u.PhotoUrl,
		HeightCm:  u.HeightCm,
		BirthDate: birthDate,
		Sex:       u.Sex,
		CreatedAt: u.CreatedAt,
	}
}

type userPage struct {
	Items  []userResponse `json:"items"`
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

func newUserPage(page types.Page[domain.User]) *userPage {
	items := make([]userResponse, len(page.Items))
	for i, u := range page.Items {
		items[i] = *userFromDomain(&u)
	}
	return &userPage{
		Items:  items,
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
}
