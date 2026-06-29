package meal_ai

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

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
	page, err := a.foodService.List(userId, foodPorts.ListParams{
		PaginationParams: types.PaginationParams{Limit: limit, Offset: 0},
		Query:            &q,
	})
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
// (base64-encoded later by pkg/openai), so the object is never made public.
//
// This works when the backend can reach the stored URL. If the R2 bucket is
// private, swap this for an S3 GetObject-based fetcher using the R2 credentials
// (see docs/ai-asistente-plan.md §4).
type httpImageFetcher struct {
	client *http.Client
}

var _ ports.ImageFetcher = (*httpImageFetcher)(nil)

func newHTTPImageFetcher() *httpImageFetcher {
	return &httpImageFetcher{client: &http.Client{Timeout: 20 * time.Second}}
}

func (f *httpImageFetcher) Fetch(ctx context.Context, url string) (ports.Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ports.Image{}, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return ports.Image{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ports.Image{}, errFetchStatus(resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ports.Image{}, err
	}
	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "image/jpeg"
	}
	return ports.Image{MimeType: mime, Data: data}, nil
}

type fetchStatusError int

func (e fetchStatusError) Error() string { return "image fetch returned non-2xx status" }

func errFetchStatus(status int) error { return fetchStatusError(status) }
