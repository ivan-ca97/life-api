package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Operation identifies the user-facing AI action that produced an interaction.
const OperationMealEstimate = "meal_estimate"

// Status values for an interaction.
const (
	StatusOK    = "ok"
	StatusError = "error"
)

// MaxInputSummaryLen caps the stored user text (no full prompt, no images).
const MaxInputSummaryLen = 500

// Interaction is one user-facing AI action (which may have fanned out into
// several provider API calls internally). Provider-agnostic so it works for any
// model vendor.
type Interaction struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	CreatedAt     time.Time
	Operation     string
	Provider      string
	Model         string
	Status        string
	ErrorType     string
	InputTokens   int64
	OutputTokens  int64
	CostUSD       float64
	LatencyMs     int
	ProviderCalls int
	CorrelationId *uuid.UUID
	InputSummary  string
	Metadata      json.RawMessage
}
