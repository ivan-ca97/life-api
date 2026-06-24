package ports

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type ImportResultItem struct {
	Date   string
	Name   string
	Status string
	Reason string
}

type ImportResult struct {
	Created  int
	Enriched int
	Skipped  int
	Blocked  int
	Results  []ImportResultItem
}

type HevyImportUseCase interface {
	Import(ctx context.Context, userId uuid.UUID, csvReader io.Reader) (*ImportResult, error)
}
