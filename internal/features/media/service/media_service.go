package service

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/media/ports"
)

var allowedContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type mediaService struct {
	storage   ports.ObjectStorage
	publicUrl string
}

var _ ports.MediaService = (*mediaService)(nil)

func NewMediaService(storage ports.ObjectStorage, publicUrl string) *mediaService {
	return &mediaService{
		storage:   storage,
		publicUrl: publicUrl,
	}
}

func (s *mediaService) GenerateUploadUrl(request ports.UploadRequest) (*ports.UploadResult, error) {
	if strings.TrimSpace(request.Filename) == "" {
		return nil, cerr.NewBadRequestError("filename is required")
	}

	ext, ok := allowedContentTypes[request.ContentType]
	if !ok {
		return nil, cerr.NewBadRequestError(fmt.Sprintf("content type %q is not allowed; use image/jpeg, image/png, or image/webp", request.ContentType))
	}

	// Prefer the extension from the original filename if it matches; otherwise use the one derived from content type
	originalExt := strings.ToLower(filepath.Ext(request.Filename))
	if originalExt != "" {
		for _, allowed := range allowedContentTypes {
			if originalExt == allowed {
				ext = originalExt
				break
			}
		}
	}

	key := fmt.Sprintf("users/%s/%s%s", request.UserId, uuid.New(), ext)

	uploadUrl, err := s.storage.GeneratePresignedPutUrl(key, request.ContentType, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	publicUrl := strings.TrimRight(s.publicUrl, "/") + "/" + key

	result := &ports.UploadResult{
		UploadUrl: uploadUrl,
		PublicUrl: publicUrl,
	}
	return result, nil
}
