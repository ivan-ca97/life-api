package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/openai"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/domain"
	aiUsagePorts "github.com/ivan-ca97/life/internal/features/ai_usage/ports"
)

// Completer is the OpenAI-backed estimator, satisfied by *openai.Client. Defined
// as a port so the use case is testable without hitting the API.
type Completer interface {
	Complete(ctx context.Context, req openai.CompletionRequest) (*openai.CompletionResult, error)
}

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

// Image is a fetched photo ready to send to the model as base64.
type Image struct {
	MimeType string
	Data     []byte
}

// ImageFetcher downloads a stored photo (e.g. from R2) so the backend can
// forward it to OpenAI without exposing the object publicly.
type ImageFetcher interface {
	Fetch(ctx context.Context, url string) (Image, error)
}

// QuotaGuard and UsageDelta are the ai_usage spend-limit contract the use case
// consumes. Aliased (not re-declared) so the ai_usage service satisfies it
// directly without a conversion adapter.
type QuotaGuard = aiUsagePorts.QuotaGuard

type UsageDelta = aiUsagePorts.UsageDelta

// Correction is a user adjustment to a prior assumption, fed back for
// re-estimation (still a single stateless call).
type Correction struct {
	Item       string
	Correction string
}

type EstimateInput struct {
	UserId            uuid.UUID
	PhotoURLs         []string
	Instructions      string
	AssumeOnlyVisible bool
	Corrections       []Correction
}

type MealEstimationUseCase interface {
	Estimate(ctx context.Context, input EstimateInput) (*domain.MealEstimate, error)
}
