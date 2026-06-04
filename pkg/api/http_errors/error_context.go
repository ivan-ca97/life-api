package http_errors

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

type contextKey string

const bagKey contextKey = "error_bag"

type errorBag struct {
	writer  http.ResponseWriter
	written bool
}

func (b *errorBag) write(err error, logger *slog.Logger) {
	if b.written {
		return
	}
	b.written = true

	statusCode := http.StatusInternalServerError
	message := "internal server error"

	var httpErr cerr.HttpError
	if errors.As(err, &httpErr) {
		statusCode = httpErr.StatusCode()
		if statusCode < 500 {
			message = httpErr.Error()
		}
	}

	level := slog.LevelError
	var customErr cerr.CustomError
	if errors.As(err, &customErr) {
		level = customErr.Severity()
	}

	attrs := []any{"error", err.Error()}
	var loggable cerr.LoggableError
	if errors.As(err, &loggable) {
		for k, v := range loggable.Log() {
			attrs = append(attrs, k, v)
		}
	}
	logger.Log(context.Background(), level, message, attrs...)

	b.writer.Header().Set("Content-Type", "application/json")
	b.writer.WriteHeader(statusCode)
	_ = json.NewEncoder(b.writer).Encode(map[string]string{"error": message})
}

// errorContextBagHandler implements HttpErrorHandler using a per-request error bag stored in context.
type errorContextBagHandler struct {
	logger *slog.Logger
}

var _ HttpErrorHandler = (*errorContextBagHandler)(nil)

func NewErrorContextBagHandler(logger *slog.Logger) *errorContextBagHandler {
	return &errorContextBagHandler{
		logger: logger,
	}
}

func (h *errorContextBagHandler) Report(r *http.Request, err error) {
	bag, ok := r.Context().Value(bagKey).(*errorBag)
	if !ok {
		return
	}
	bag.write(err, h.logger)
}

// Middleware initializes the error bag in the request context.
// Must wrap all routes that use HttpErrorHandler.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bag := &errorBag{
			writer: w,
		}
		ctx := context.WithValue(r.Context(), bagKey, bag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
