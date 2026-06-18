package handler

import (
	"time"

	"github.com/ivan-ca97/life/internal/features/measurements/domain"
)

type measurementResponse struct {
	Date      string    `json:"date"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Notes     string    `json:"notes"`
	UpdatedAt time.Time `json:"updated_at"`
}

type measurementListResponse struct {
	Items []measurementResponse `json:"items"`
}

func measurementFromDomain(m *domain.BodyMeasurement) *measurementResponse {
	return &measurementResponse{
		Date:      m.Date.Format("2006-01-02"),
		Type:      m.Type,
		Value:     m.Value,
		Notes:     m.Notes,
		UpdatedAt: m.UpdatedAt,
	}
}
