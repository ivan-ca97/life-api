package ports

import "github.com/google/uuid"

// WeightLookup retrieves the most recent weight for a user, used to compute
// calories burned from steps (steps × 0.0005 × weight_kg).
type WeightLookup interface {
	LatestWeightKg(userId uuid.UUID) (float64, bool, error)
}
