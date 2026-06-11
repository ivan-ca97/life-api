package media

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (f *mediaFeature) ProtectedRoutes(r chi.Router) {
	r.Post("/media/upload-url", endpoint.JSON(f.errorHandler, f.mediaHandler.GenerateUploadURL))
}
