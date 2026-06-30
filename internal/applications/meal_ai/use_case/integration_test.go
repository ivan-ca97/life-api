package use_case

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/openai"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
)

// localImageFetcher reads a photo from disk so the integration test does not need
// a public URL. MIME type is inferred from the file extension.
type localImageFetcher struct{ path string }

func (f localImageFetcher) Fetch(ctx context.Context, _ string) (ports.Image, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return ports.Image{}, err
	}
	mime := "image/jpeg"
	switch {
	case strings.HasSuffix(strings.ToLower(f.path), ".png"):
		mime = "image/png"
	case strings.HasSuffix(strings.ToLower(f.path), ".webp"):
		mime = "image/webp"
	}
	return ports.Image{MimeType: mime, Data: data}, nil
}

// fixedFoodSearch returns a tiny fake catalog so we can watch the model match
// against it (and suggest new foods for anything missing).
type fixedFoodSearch struct{}

func (fixedFoodSearch) Search(userId uuid.UUID, query string, limit int) ([]ports.FoodCandidate, error) {
	cal, p, c, fat := 165.0, 31.0, 0.0, 3.6
	riceCal, riceP, riceC, riceFat := 130.0, 2.7, 28.0, 0.3
	return []ports.FoodCandidate{
		{Id: uuid.New(), Name: "Pechuga de pollo a la plancha", DefaultCalories: &cal, DefaultProteinGrams: &p, DefaultCarbsGrams: &c, DefaultFatGrams: &fat, MeasurementType: "mass", BaseQuantity: 100, BaseUnit: "g"},
		{Id: uuid.New(), Name: "Arroz blanco cocido", DefaultCalories: &riceCal, DefaultProteinGrams: &riceP, DefaultCarbsGrams: &riceC, DefaultFatGrams: &riceFat, MeasurementType: "mass", BaseQuantity: 100, BaseUnit: "g"},
	}, nil
}

// TestEstimate_Integration hits the real OpenAI API. It is skipped unless
// OPENAI_API_KEY is set, so it never runs in normal `go test`.
//
// Run it explicitly (PowerShell):
//
//	$env:OPENAI_API_KEY = "sk-..."
//	$env:OPENAI_TEST_IMAGE = "C:\ruta\a\una\foto.jpg"   # opcional; sin esto, prueba solo texto
//	go test -run Integration -v ./internal/applications/meal_ai/use_case/
func TestEstimate_Integration(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("set OPENAI_API_KEY to run the OpenAI integration test")
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o"
	}

	client := openai.NewClient(openai.Config{APIKey: apiKey, Model: model, Timeout: 90 * time.Second})

	input := ports.EstimateInput{UserId: uuid.New(), AssumeOnlyVisible: true}
	imagePath := os.Getenv("OPENAI_TEST_IMAGE")
	var fetcher ports.ImageFetcher = localImageFetcher{path: imagePath}
	if imagePath != "" {
		input.PhotoURLs = []string{"file://local"} // value unused; fetcher reads the path
	} else {
		// No photo: drive it with a text description instead.
		input.Instructions = "Comí unos 200g de pechuga de pollo a la plancha, 150g de arroz blanco y un poco de pollo frito."
	}

	uc := NewMealEstimationUseCase(client, fixedFoodSearch{}, fetcher, noopQuota{}, stubAuthorizer{}, model)

	estimate, err := uc.Estimate(context.Background(), input)
	if err != nil {
		if cause := errors.Unwrap(err); cause != nil {
			t.Fatalf("real estimation failed: %v\n  cause: %v", err, cause)
		}
		t.Fatalf("real estimation failed: %v", err)
	}

	t.Logf("matched items: %d, new-food suggestions: %d", len(estimate.MatchedItems), len(estimate.NewFoodSuggestions))
	for _, m := range estimate.MatchedItems {
		t.Logf("  MATCH %s ~%.0f%s (conf %s) warnings=%v", m.FoodName, m.EstimatedQuantity, m.Unit, m.Confidence, m.SanityWarnings)
	}
	for _, s := range estimate.NewFoodSuggestions {
		t.Logf("  NEW   %s ~%.0f%s (conf %s)", s.Name, s.EstimatedQuantity, s.Unit, s.Confidence)
	}
	t.Logf("totals: %.0f kcal / P%.0f C%.0f F%.0f", estimate.Totals.Calories, estimate.Totals.ProteinGrams, estimate.Totals.CarbsGrams, estimate.Totals.FatGrams)
	t.Logf("assumptions: %v", estimate.Assumptions)
	t.Logf("usage: %d in + %d out tokens = $%.4f (model %s)", estimate.Usage.InputTokens, estimate.Usage.OutputTokens, estimate.Usage.CostUSD, estimate.Usage.Model)

	if len(estimate.MatchedItems) == 0 && len(estimate.NewFoodSuggestions) == 0 {
		t.Error("expected at least one matched item or new-food suggestion")
	}
}

type noopQuota struct{}

func (noopQuota) CheckQuota(userId uuid.UUID) error                          { return nil }
func (noopQuota) RecordUsage(userId uuid.UUID, delta ports.UsageDelta) error { return nil }
