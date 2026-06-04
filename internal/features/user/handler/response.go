package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
)

type userResponse struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

func userFromDomain(u *domain.User) *userResponse {
	return &userResponse{
		Id:        u.Id,
		Email:     u.Email,
		Active:    u.Active,
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
