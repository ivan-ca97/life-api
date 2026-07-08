package llm

import (
	"context"
	"encoding/json"
)

// Tool is something the model can invoke during a completion.
type Tool interface {
	Name() string
	Description() string
	Schema() json.RawMessage
	Call(ctx context.Context, arguments json.RawMessage) (string, error)
}
