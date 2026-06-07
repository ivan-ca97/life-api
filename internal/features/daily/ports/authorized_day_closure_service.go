package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthorizedDayClosureService interface {
	Close(ctx context.Context, ownerId uuid.UUID, date time.Time) error
	Open(ctx context.Context, ownerId uuid.UUID, date time.Time) error
	IsClosed(ctx context.Context, ownerId uuid.UUID, date time.Time) (bool, error)
}
