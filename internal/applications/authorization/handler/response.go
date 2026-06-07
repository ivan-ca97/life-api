package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/authorization/domain"
)

type shareResponse struct {
	Id           uuid.UUID `json:"id"`
	OwnerId      uuid.UUID `json:"owner_id"`
	GranteeId    uuid.UUID `json:"grantee_id"`
	ResourceType string    `json:"resource_type"`
	CanWrite     bool      `json:"can_write"`
	CreatedAt    time.Time `json:"created_at"`
}

type shareListResponse struct {
	Items []shareResponse `json:"items"`
}

func shareFromDomain(s *domain.Share) *shareResponse {
	return &shareResponse{
		Id:           s.Id,
		OwnerId:      s.OwnerId,
		GranteeId:    s.GranteeId,
		ResourceType: s.ResourceType,
		CanWrite:     s.CanWrite,
		CreatedAt:    s.CreatedAt,
	}
}

func shareListFromDomain(shares []domain.Share) *shareListResponse {
	items := make([]shareResponse, len(shares))
	for i, s := range shares {
		items[i] = *shareFromDomain(&s)
	}
	return &shareListResponse{
		Items: items,
	}
}
