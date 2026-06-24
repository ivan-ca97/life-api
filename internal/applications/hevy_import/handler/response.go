package handler

import (
	"github.com/ivan-ca97/life/internal/applications/hevy_import/ports"
)

type importResultItemResponse struct {
	Date   string `json:"date"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type importResponse struct {
	Created  int                        `json:"created"`
	Enriched int                        `json:"enriched"`
	Skipped  int                        `json:"skipped"`
	Blocked  int                        `json:"blocked"`
	Results  []importResultItemResponse `json:"results"`
}

func importResponseFromResult(result *ports.ImportResult) *importResponse {
	items := make([]importResultItemResponse, len(result.Results))
	for i, r := range result.Results {
		items[i] = importResultItemResponse{
			Date:   r.Date,
			Name:   r.Name,
			Status: r.Status,
			Reason: r.Reason,
		}
	}
	return &importResponse{
		Created:  result.Created,
		Enriched: result.Enriched,
		Skipped:  result.Skipped,
		Blocked:  result.Blocked,
		Results:  items,
	}
}
