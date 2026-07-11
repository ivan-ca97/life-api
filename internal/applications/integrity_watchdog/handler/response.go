package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/ports"
	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/scheduler"
)

type crossContextPhotoResponse struct {
	PhotoId    uuid.UUID `json:"photo_id"`
	MealId     uuid.UUID `json:"meal_id"`
	MealItemId uuid.UUID `json:"meal_item_id"`
}

type itemGroupResponse struct {
	MealId     uuid.UUID `json:"meal_id"`
	MealItemId uuid.UUID `json:"meal_item_id"`
}

type invalidFoodUnitResponse struct {
	FoodId   uuid.UUID `json:"food_id"`
	BaseUnit string    `json:"base_unit"`
}

type dbIntegrityResponse struct {
	Clean                    bool                        `json:"clean"`
	CrossContextPhotos       []crossContextPhotoResponse `json:"cross_context_photos"`
	MealGroupsMissingPrimary []uuid.UUID                 `json:"meal_groups_missing_primary"`
	ItemGroupsMissingPrimary []itemGroupResponse         `json:"item_groups_missing_primary"`
	InvalidFoodBaseUnits     []invalidFoodUnitResponse   `json:"invalid_food_base_units"`
}

type r2OrphanResponse struct {
	Clean       bool     `json:"clean"`
	OrphanKeys  []string `json:"orphan_keys"`
	DeletedKeys []string `json:"deleted_keys"`
	BrokenRefs  []string `json:"broken_refs"`
}

type lastRunResponse struct {
	StartedAt  time.Time            `json:"started_at"`
	FinishedAt time.Time            `json:"finished_at"`
	DurationMs int64                `json:"duration_ms"`
	Error      *string              `json:"error,omitempty"`
	DB         *dbIntegrityResponse `json:"db,omitempty"`
	R2         *r2OrphanResponse    `json:"r2,omitempty"`
}

type statusResponse struct {
	Running         bool             `json:"running"`
	IntervalSeconds int64            `json:"interval_seconds"`
	LastRun         *lastRunResponse `json:"last_run,omitempty"`
}

func buildStatusResponse(s *scheduler.Scheduler) *statusResponse {
	response := &statusResponse{
		Running:         s.IsRunning(),
		IntervalSeconds: int64(s.Period().Seconds()),
	}
	last := s.LastResult()
	if last != nil {
		run := &lastRunResponse{
			StartedAt:  last.StartedAt,
			FinishedAt: last.FinishedAt,
			DurationMs: last.FinishedAt.Sub(last.StartedAt).Milliseconds(),
		}
		if last.Err != nil {
			message := last.Err.Error()
			run.Error = &message
		}
		if last.DB != nil {
			run.DB = buildDBResponse(last.DB.CrossContextPhotos, last.DB.MealGroupsMissingPrimary, last.DB.ItemGroupsMissingPrimary, last.DB.InvalidFoodBaseUnits)
		}
		if last.R2 != nil {
			run.R2 = &r2OrphanResponse{
				Clean:       last.R2.IsClean(),
				OrphanKeys:  nullSlice(last.R2.OrphanKeys),
				DeletedKeys: nullSlice(last.R2.DeletedKeys),
				BrokenRefs:  nullSlice(last.R2.BrokenRefs),
			}
		}
		response.LastRun = run
	}
	return response
}

func buildDBResponse(cross []ports.CrossContextPhoto, mealMissing []uuid.UUID, itemMissing []ports.ItemGroup, invalidUnits []ports.InvalidFoodUnit) *dbIntegrityResponse {
	crossResp := make([]crossContextPhotoResponse, len(cross))
	for i, c := range cross {
		crossResp[i] = crossContextPhotoResponse{PhotoId: c.PhotoId, MealId: c.MealId, MealItemId: c.MealItemId}
	}
	itemResp := make([]itemGroupResponse, len(itemMissing))
	for i, g := range itemMissing {
		itemResp[i] = itemGroupResponse{MealId: g.MealId, MealItemId: g.MealItemId}
	}
	unitResp := make([]invalidFoodUnitResponse, len(invalidUnits))
	for i, u := range invalidUnits {
		unitResp[i] = invalidFoodUnitResponse{FoodId: u.FoodId, BaseUnit: u.BaseUnit}
	}
	clean := len(cross) == 0 && len(mealMissing) == 0 && len(itemMissing) == 0 && len(invalidUnits) == 0
	return &dbIntegrityResponse{
		Clean:                    clean,
		CrossContextPhotos:       crossResp,
		MealGroupsMissingPrimary: nullSlice(mealMissing),
		ItemGroupsMissingPrimary: itemResp,
		InvalidFoodBaseUnits:     unitResp,
	}
}

func nullSlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
