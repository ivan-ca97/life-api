package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
)

type goalResponse struct {
	Id                   uuid.UUID `json:"id"`
	DailyCalories        *float64  `json:"daily_calories,omitempty"`
	DailyProteinGrams    *float64  `json:"daily_protein_grams,omitempty"`
	DailyCarbsGrams      *float64  `json:"daily_carbs_grams,omitempty"`
	DailyFatGrams        *float64  `json:"daily_fat_grams,omitempty"`
	DailyFiberGrams      *float64  `json:"daily_fiber_grams,omitempty"`
	DailySteps           *int      `json:"daily_steps,omitempty"`
	DailyExerciseMinutes *int      `json:"daily_exercise_minutes,omitempty"`
	TargetWeightKg       *float64  `json:"target_weight_kg,omitempty"`
	StartedAt            time.Time `json:"started_at"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func goalFromDomain(g *domain.Goal) *goalResponse {
	return &goalResponse{
		Id:                   g.Id,
		DailyCalories:        g.DailyCalories,
		DailyProteinGrams:    g.DailyProteinGrams,
		DailyCarbsGrams:      g.DailyCarbsGrams,
		DailyFatGrams:        g.DailyFatGrams,
		DailyFiberGrams:      g.DailyFiberGrams,
		DailySteps:           g.DailySteps,
		DailyExerciseMinutes: g.DailyExerciseMinutes,
		TargetWeightKg:       g.TargetWeightKg,
		StartedAt:            g.StartedAt,
		CreatedAt:            g.CreatedAt,
		UpdatedAt:            g.UpdatedAt,
	}
}

type goalMetricResponse struct {
	Target      float64 `json:"target"`
	Average     float64 `json:"average"`
	DaysMet     int     `json:"days_met"`
	DaysTracked int     `json:"days_tracked"`
	DaysTotal   int     `json:"days_total"`
}

type weightProgressResponse struct {
	TargetKg  float64  `json:"target_kg"`
	CurrentKg *float64 `json:"current_kg,omitempty"`
}

type goalProgressResponse struct {
	From                 string                  `json:"from"`
	To                   string                  `json:"to"`
	DaysTotal            int                     `json:"days_total"`
	Goal                 *goalResponse           `json:"goal"`
	DailyCalories        *goalMetricResponse     `json:"daily_calories,omitempty"`
	DailyProteinGrams    *goalMetricResponse     `json:"daily_protein_grams,omitempty"`
	DailyCarbsGrams      *goalMetricResponse     `json:"daily_carbs_grams,omitempty"`
	DailyFatGrams        *goalMetricResponse     `json:"daily_fat_grams,omitempty"`
	DailyFiberGrams      *goalMetricResponse     `json:"daily_fiber_grams,omitempty"`
	DailySteps           *goalMetricResponse     `json:"daily_steps,omitempty"`
	DailyExerciseMinutes *goalMetricResponse     `json:"daily_exercise_minutes,omitempty"`
	WeightProgress       *weightProgressResponse `json:"weight_progress,omitempty"`
}

func metricFromDomain(m *domain.GoalMetric) *goalMetricResponse {
	if m == nil {
		return nil
	}
	return &goalMetricResponse{
		Target:      m.Target,
		Average:     m.Average,
		DaysMet:     m.DaysMet,
		DaysTracked: m.DaysTracked,
		DaysTotal:   m.DaysTotal,
	}
}

func goalProgressFromDomain(p *domain.GoalProgress) *goalProgressResponse {
	resp := &goalProgressResponse{
		From:                 p.From.Format("2006-01-02"),
		To:                   p.To.Format("2006-01-02"),
		DaysTotal:            p.DaysTotal,
		Goal:                 goalFromDomain(p.Goal),
		DailyCalories:        metricFromDomain(p.DailyCalories),
		DailyProteinGrams:    metricFromDomain(p.DailyProteinGrams),
		DailyCarbsGrams:      metricFromDomain(p.DailyCarbsGrams),
		DailyFatGrams:        metricFromDomain(p.DailyFatGrams),
		DailyFiberGrams:      metricFromDomain(p.DailyFiberGrams),
		DailySteps:           metricFromDomain(p.DailySteps),
		DailyExerciseMinutes: metricFromDomain(p.DailyExerciseMinutes),
	}
	if p.WeightProgress != nil {
		resp.WeightProgress = &weightProgressResponse{
			TargetKg:  p.WeightProgress.TargetKg,
			CurrentKg: p.WeightProgress.CurrentKg,
		}
	}
	return resp
}
