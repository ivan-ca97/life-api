package llm

import (
	"errors"
	"fmt"
)

// Completion problems reported via finish reason (the call succeeded, but the
// answer is unusable). Adapters map their provider's reasons to these.
var (
	ErrTruncated       = errors.New("llm: response truncated before completion")
	ErrContentFiltered = errors.New("llm: response blocked by the content filter")
)

type ProviderError struct {
	Provider   string
	Type       string
	Message    string
	StatusCode int
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("llm: %s error (status %d, type %q): %s", e.Provider, e.StatusCode, e.Type, e.Message)
}
