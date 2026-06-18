package handler

import (
	"time"

	"github.com/ivan-ca97/life/internal/features/steps/domain"
)

type stepsResponse struct {
	Date           string   `json:"date"`
	Steps          int      `json:"steps"`
	Source         string   `json:"source"`
	CaloriesBurned *float64 `json:"calories_burned,omitempty"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func stepsFromDomain(s *domain.DailySteps, weightKg *float64) *stepsResponse {
	resp := &stepsResponse{
		Date:      s.Date.Format("2006-01-02"),
		Steps:     s.Steps,
		Source:    s.Source,
		UpdatedAt: s.UpdatedAt,
	}
	if weightKg != nil {
		cal := float64(s.Steps) * 0.0005 * *weightKg
		resp.CaloriesBurned = &cal
	}
	return resp
}

type stepsListResponse struct {
	Items []stepsResponse `json:"items"`
}
