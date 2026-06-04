package domain

import "time"

type Correction struct {
	Date            time.Time
	Calories        *float64
	ProteinGrams    *float64
	CarbsGrams      *float64
	FatGrams        *float64
	FiberGrams      *float64
	CaloriesBurned  *float64
	Steps           *int
	DurationSeconds *int
	DistanceMeters  *float64
	Notes           string
}

func (c *Correction) HasMealFields() bool {
	return c.Calories != nil || c.ProteinGrams != nil || c.CarbsGrams != nil ||
		c.FatGrams != nil || c.FiberGrams != nil
}

func (c *Correction) HasExerciseFields() bool {
	return c.CaloriesBurned != nil || c.Steps != nil ||
		c.DurationSeconds != nil || c.DistanceMeters != nil
}
