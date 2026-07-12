package checks

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/ports"
)

type DBIntegrityResult struct {
	CrossContextPhotos       []ports.CrossContextPhoto
	MealGroupsMissingPrimary []uuid.UUID
	ItemGroupsMissingPrimary []ports.ItemGroup
	InvalidFoodBaseUnits     []ports.InvalidFoodUnit
}

func (r *DBIntegrityResult) IsClean() bool {
	return len(r.CrossContextPhotos) == 0 &&
		len(r.MealGroupsMissingPrimary) == 0 &&
		len(r.ItemGroupsMissingPrimary) == 0 &&
		len(r.InvalidFoodBaseUnits) == 0
}

type DBIntegrityCheck struct {
	repository ports.WatchdogRepository
}

func NewDBIntegrityCheck(repository ports.WatchdogRepository) *DBIntegrityCheck {
	return &DBIntegrityCheck{repository: repository}
}

func (c *DBIntegrityCheck) Run() (*DBIntegrityResult, error) {
	cross, err := c.repository.CrossContextPhotos()
	if err != nil {
		return nil, fmt.Errorf("cross-context photos: %w", err)
	}

	mealMissing, err := c.repository.MealGroupsMissingPrimary()
	if err != nil {
		return nil, fmt.Errorf("meal groups missing primary: %w", err)
	}

	itemMissing, err := c.repository.ItemGroupsMissingPrimary()
	if err != nil {
		return nil, fmt.Errorf("item groups missing primary: %w", err)
	}

	invalidUnits, err := c.repository.InvalidFoodBaseUnits()
	if err != nil {
		return nil, fmt.Errorf("invalid food base units: %w", err)
	}

	result := &DBIntegrityResult{
		CrossContextPhotos:       cross,
		MealGroupsMissingPrimary: mealMissing,
		ItemGroupsMissingPrimary: itemMissing,
		InvalidFoodBaseUnits:     invalidUnits,
	}
	return result, nil
}
