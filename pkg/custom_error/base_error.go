package custom_error

import (
	"log/slog"
	"net/http"
)

// BaseHttpError is the concrete base type for all domain errors.
// Domain packages create singletons via the constructors below.
type BaseHttpError struct {
	msg        string
	statusCode int
	severity   slog.Level
}

func (e *BaseHttpError) Error() string        { return e.msg }
func (e *BaseHttpError) StatusCode() int      { return e.statusCode }
func (e *BaseHttpError) Severity() slog.Level { return e.severity }

func NewNotFoundError(entity string) *BaseHttpError {
	return &BaseHttpError{
		msg:        entity + " not found",
		statusCode: http.StatusNotFound,
		severity:   slog.LevelWarn,
	}
}

func NewConflictError(msg string) *BaseHttpError {
	return &BaseHttpError{
		msg:        msg,
		statusCode: http.StatusConflict,
		severity:   slog.LevelWarn,
	}
}

func NewUnauthorizedError(msg string) *BaseHttpError {
	return &BaseHttpError{
		msg:        msg,
		statusCode: http.StatusUnauthorized,
		severity:   slog.LevelWarn,
	}
}

func NewForbiddenError(msg string) *BaseHttpError {
	return &BaseHttpError{
		msg:        msg,
		statusCode: http.StatusForbidden,
		severity:   slog.LevelWarn,
	}
}

func NewBadRequestError(msg string) *BaseHttpError {
	return &BaseHttpError{
		msg:        msg,
		statusCode: http.StatusBadRequest,
		severity:   slog.LevelWarn,
	}
}

func NewTooManyRequestsError(msg string) *BaseHttpError {
	return &BaseHttpError{
		msg:        msg,
		statusCode: http.StatusTooManyRequests,
		severity:   slog.LevelWarn,
	}
}

// internalError wraps an unexpected error with context for logging.
// The public message is always "internal server error" — the cause is only logged.
type internalError struct {
	*BaseHttpError
	cause      error
	logContext string
}

func (e *internalError) Log() map[string]string {
	m := map[string]string{"context": e.logContext}
	if e.cause != nil {
		m["cause"] = e.cause.Error()
	}
	return m
}

// Unwrap exposes the underlying cause for errors.Is/As and logging, while the
// public Error() message stays the generic "internal server error".
func (e *internalError) Unwrap() error { return e.cause }

func NewInternalError(logContext string, cause error) *internalError {
	return &internalError{
		BaseHttpError: &BaseHttpError{
			msg:        "internal server error",
			statusCode: http.StatusInternalServerError,
			severity:   slog.LevelError,
		},
		logContext: logContext,
		cause:      cause,
	}
}
