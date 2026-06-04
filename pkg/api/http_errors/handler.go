package http_errors

import "net/http"

type HttpErrorHandler interface {
	Report(r *http.Request, err error)
}
