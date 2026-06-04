package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
)

type weightEntryResponse struct {
	Id                uuid.UUID `json:"id"`
	Date              string    `json:"date"`
	WeightKg          float64   `json:"weight_kg"`
	BodyFatPercentage *float64  `json:"body_fat_percentage,omitempty"`
	Notes             string    `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func weightEntryFromDomain(e *domain.WeightEntry) *weightEntryResponse {
	return &weightEntryResponse{
		Id:                e.Id,
		Date:              e.Date.Format("2006-01-02"),
		WeightKg:          e.WeightKg,
		BodyFatPercentage: e.BodyFatPercentage,
		Notes:             e.Notes,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

type weightEntryPage struct {
	Items  []weightEntryResponse `json:"items"`
	Total  int64                 `json:"total"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}

func newWeightEntryPage(page types.Page[domain.WeightEntry]) *weightEntryPage {
	items := make([]weightEntryResponse, len(page.Items))
	for i, e := range page.Items {
		items[i] = *weightEntryFromDomain(&e)
	}
	return &weightEntryPage{
		Items:  items,
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
}
