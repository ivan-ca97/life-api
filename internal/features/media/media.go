package media

import (
	"github.com/ivan-ca97/life/pkg/api/http_errors"

	"github.com/ivan-ca97/life/internal/features/media/handler"
	"github.com/ivan-ca97/life/internal/features/media/service"
	"github.com/ivan-ca97/life/internal/features/media/storage"
)

type mediaFeature struct {
	mediaHandler handler.MediaHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewMediaFeature(
	accountId string,
	accessKeyId string,
	secretAccessKey string,
	bucket string,
	publicURL string,
	errorHandler http_errors.HttpErrorHandler,
) *mediaFeature {
	r2 := storage.NewR2Storage(accountId, accessKeyId, secretAccessKey, bucket)
	mediaService := service.NewMediaService(r2, publicURL)
	mediaHandler := handler.NewMediaHandler(mediaService)

	return &mediaFeature{
		mediaHandler: mediaHandler,
		errorHandler: errorHandler,
	}
}
