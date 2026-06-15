package handler

import "github.com/ivan-ca97/life/internal/applications/fitness_advisor/domain"

type estimateResponse struct {
	Type              string  `json:"type"`
	Value             float64 `json:"value"`
	EstimatedCalories float64 `json:"estimated_calories"`
	WeightKg          float64 `json:"weight_kg"`
}

func estimateResponseFromDomain(result *domain.EstimateResult) *estimateResponse {
	return &estimateResponse{
		Type:              string(result.Type),
		Value:             result.Value,
		EstimatedCalories: result.EstimatedCalories,
		WeightKg:          result.WeightKg,
	}
}
