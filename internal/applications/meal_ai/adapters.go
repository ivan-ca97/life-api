package meal_ai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"

	nutritionPorts "github.com/ivan-ca97/life/internal/features/nutrition/ports"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
)

// foodSearchAdapter implements ports.FoodSearch over the food feature's service.
type foodSearchAdapter struct {
	foodService nutritionPorts.FoodService
}

var _ ports.FoodSearch = (*foodSearchAdapter)(nil)

// Search runs each keyword as a separate name match and returns the
// deduplicated union (by food id), capped at limit. Searching "lomo" and
// "vaca" separately both surface "Lomo de vaca"; it is returned only once.
func (a *foodSearchAdapter) Search(userId uuid.UUID, queries []string, limit int) ([]ports.FoodCandidate, error) {
	seen := make(map[uuid.UUID]bool)
	candidates := make([]ports.FoodCandidate, 0, limit)
	for _, query := range queries {
		q := strings.TrimSpace(query)
		if q == "" {
			continue
		}
		params := nutritionPorts.ListFoodParams{
			PaginationParams: types.PaginationParams{
				Limit:  limit,
				Offset: 0,
			},
			Query: &q,
		}
		page, err := a.foodService.List(userId, params)
		if err != nil {
			return nil, err
		}
		for _, f := range page.Items {
			if seen[f.Id] {
				continue
			}
			seen[f.Id] = true
			candidate := ports.FoodCandidate{
				Id:                  f.Id,
				Name:                f.Name,
				DefaultCalories:     f.DefaultCalories,
				DefaultProteinGrams: f.DefaultProteinGrams,
				DefaultCarbsGrams:   f.DefaultCarbsGrams,
				DefaultFatGrams:     f.DefaultFatGrams,
				DefaultFiberGrams:   f.DefaultFiberGrams,
				MeasurementType:     string(f.Dimension()),
				BaseQuantity:        f.BaseQuantity,
				BaseUnit:            f.BaseUnit.String(),
			}
			candidates = append(candidates, candidate)
			if len(candidates) >= limit {
				return candidates, nil
			}
		}
	}
	return candidates, nil
}

// httpImageFetcher downloads a stored photo over HTTP and forwards the raw bytes
// (base64-encoded later by the provider adapter), so the object is never made
// public.
//
// This works when the backend can reach the stored URL. If the R2 bucket is
// private, swap this for an S3 GetObject-based fetcher using the R2 credentials
// (see docs/ai-asistente-plan.md §4).
type httpImageFetcher struct {
	client *http.Client
}

var _ ports.ImageFetcher = (*httpImageFetcher)(nil)

func newHTTPImageFetcher() *httpImageFetcher {
	return &httpImageFetcher{client: &http.Client{
		Timeout: 20 * time.Second,
	}}
}

func (f *httpImageFetcher) Fetch(ctx context.Context, url string) (llm.Image, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return llm.Image{}, err
	}
	response, err := f.client.Do(request)
	if err != nil {
		return llm.Image{}, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return llm.Image{}, errFetchStatus(response.StatusCode)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return llm.Image{}, err
	}
	mediaType := response.Header.Get("Content-Type")
	if mediaType == "" {
		mediaType = "image/jpeg"
	}
	result := llm.Image{
		MediaType: mediaType,
		Data:      data,
	}
	return result, nil
}

type fetchStatusError int

func (e fetchStatusError) Error() string { return "image fetch returned non-2xx status" }

func errFetchStatus(status int) error { return fetchStatusError(status) }
