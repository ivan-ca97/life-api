package domain

import "time"

// ModelPrice is a provider/model's token price, valid from EffectiveFrom onward
// (until a later row supersedes it).
type ModelPrice struct {
	Provider         string
	Model            string
	InputPerMillion  float64
	OutputPerMillion float64
	EffectiveFrom    time.Time
}
