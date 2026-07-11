package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/domain"
	aiUsagePorts "github.com/ivan-ca97/life/internal/features/ai_usage/ports"
)

// FoodCandidate is a catalog entry returned to the model so it can match a
// detected food and sanity-check the stored macros.
type FoodCandidate struct {
	Id                  uuid.UUID
	Name                string
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	MeasurementType     string
	BaseQuantity        float64
	BaseUnit            string
}

// FoodSearch looks up a user's food catalog. Satisfied by an adapter over the
// food feature's service.
type FoodSearch interface {
	Search(userId uuid.UUID, query string, limit int) ([]FoodCandidate, error)
}

// ImageFetcher downloads a stored photo (e.g. from R2) and returns it ready to
// send to the model, so the object is never made public.
type ImageFetcher interface {
	Fetch(ctx context.Context, url string) (llm.Image, error)
}

// QuotaGuard, UsageDelta, InteractionLogger and InteractionEntry are the ai_usage
// contracts the use case consumes, aliased so the ai_usage service satisfies them
// directly.
type QuotaGuard = aiUsagePorts.QuotaGuard

type UsageDelta = aiUsagePorts.UsageDelta

type InteractionLogger = aiUsagePorts.InteractionLogger

type InteractionEntry = aiUsagePorts.InteractionEntry

// Pricer converts token usage to a USD cost. Defined here (the consumer);
// satisfied by the ai_usage service.
type Pricer interface {
	CostUSD(provider, model string, inputTokens, outputTokens int64, at time.Time) (float64, error)
}

// Correction is a user adjustment to a prior assumption, fed back for
// re-estimation (still a single stateless call).
type Correction struct {
	Item       string
	Correction string
}

type EstimateInput struct {
	UserId            uuid.UUID
	PhotoUrls         []string
	Instructions      string
	AssumeOnlyVisible bool
	Corrections       []Correction
}

type MealEstimationUseCase interface {
	Estimate(ctx context.Context, input EstimateInput) (*domain.MealEstimate, error)
}
