package endpoint

import (
	"encoding/json"
	"net/http"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
)

type JSONHandler[T any] func(*http.Request) (*T, int, error)

func JSON[T any](errorHandler http_errors.HttpErrorHandler, handler JSONHandler[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, status, err := handler(r)
		if err != nil {
			errorHandler.Report(r, err)
			return
		}
		err = writeJSON(w, status, result)
		if err != nil {
			errorHandler.Report(r, err)
		}
	}
}

func writeJSON[T any](w http.ResponseWriter, status int, data *T) error {
	if data == nil {
		w.WriteHeader(status)
		return nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}
