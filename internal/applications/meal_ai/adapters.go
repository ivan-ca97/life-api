package meal_ai

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"

	foodPorts "github.com/ivan-ca97/life/internal/features/food/ports"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
)

// foodSearchAdapter implements ports.FoodSearch over the food feature's service.
type foodSearchAdapter struct {
	foodService foodPorts.FoodService
}

var _ ports.FoodSearch = (*foodSearchAdapter)(nil)

func (a *foodSearchAdapter) Search(userId uuid.UUID, query string, limit int) ([]ports.FoodCandidate, error) {
	q := query
	params := foodPorts.ListParams{
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
	candidates := make([]ports.FoodCandidate, len(page.Items))
	for i, f := range page.Items {
		candidates[i] = ports.FoodCandidate{
			Id:                  f.Id,
			Name:                f.Name,
			DefaultCalories:     f.DefaultCalories,
			DefaultProteinGrams: f.DefaultProteinGrams,
			DefaultCarbsGrams:   f.DefaultCarbsGrams,
			DefaultFatGrams:     f.DefaultFatGrams,
			DefaultFiberGrams:   f.DefaultFiberGrams,
			MeasurementType:     f.MeasurementType,
			BaseQuantity:        f.BaseQuantity,
			BaseUnit:            f.BaseUnit,
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
