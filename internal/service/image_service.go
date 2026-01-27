package service

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
	apperrors "go-upload/internal/domain/errors"
	"go-upload/internal/dto"
	"go-upload/internal/repository"
	"gorm.io/gorm"
)

const (
	DefaultQuality = 80
	MaxQuality     = 100
)

var supportedFormats = map[string]bimg.ImageType{
	"webp": bimg.WEBP,
	"jpeg": bimg.JPEG,
	"jpg":  bimg.JPEG,
	"png":  bimg.PNG,
	"avif": bimg.AVIF,
}

var formatToMime = map[string]string{
	"webp": "image/webp",
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"avif": "image/avif",
}

var mimeTypeMap = map[string]bimg.ImageType{
	"image/jpeg": bimg.JPEG,
	"image/png":  bimg.PNG,
	"image/webp": bimg.WEBP,
}

// imageService handles image processing operations
type imageService struct {
	uploadRepo repository.UploadRepository
}

// NewImageService creates a new image service
func NewImageService(uploadRepo repository.UploadRepository) ImageService {
	return &imageService{
		uploadRepo: uploadRepo,
	}
}

// ProcessImage retrieves and processes an image with specified parameters
func (s *imageService) ProcessImage(ctx context.Context, fileID uuid.UUID, params dto.ImageParams) ([]byte, string, error) {
	// Find upload record
	upload, err := s.uploadRepo.FindByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperrors.ErrFileNotFound
		}
		return nil, "", err
	}

	// Detect MIME type
	mimeType := mime.TypeByExtension(filepath.Ext(upload.FilePath))
	if !strings.HasPrefix(mimeType, "image/") {
		return nil, "", apperrors.ErrInvalidFileType
	}

	// Process image
	processedImage, outputMimeType, err := s.processImageWithOptions(upload.FilePath, params, mimeType)
	if err != nil {
		return nil, "", err
	}

	return processedImage, outputMimeType, nil
}

// ValidateImageParams validates image processing parameters
func (s *imageService) ValidateImageParams(params *dto.ImageParams) error {
	// Set default quality if not specified
	if params.Quality == 0 {
		params.Quality = DefaultQuality
	}

	// Validate quality
	if params.Quality < 1 || params.Quality > MaxQuality {
		return apperrors.ErrInvalidImageParams
	}

	// Validate width
	if params.Width < 0 {
		return apperrors.ErrInvalidImageParams
	}

	// Validate format if specified
	if params.Format != "" && params.Format != "original" {
		if _, supported := supportedFormats[params.Format]; !supported {
			return apperrors.ErrInvalidImageFormat
		}
	}

	return nil
}

func (s *imageService) processImageWithOptions(filePath string, params dto.ImageParams, contentType string) ([]byte, string, error) {
	// Read image file
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	// Create bimg image
	img := bimg.NewImage(imageData)
	if img == nil {
		return nil, "", apperrors.ErrImageProcessingFailed
	}

	// Set processing options
	options := bimg.Options{
		Quality: params.Quality,
	}

	// Set width if specified
	if params.Width > 0 {
		options.Width = params.Width
		options.Enlarge = false
		options.Force = false
	}

	// Set format if specified
	outputMimeType := contentType
	if params.Format != "" && params.Format != "original" {
		if imageType, supported := supportedFormats[params.Format]; supported {
			options.Type = imageType
			outputMimeType = formatToMime[params.Format]
		} else {
			return nil, "", apperrors.ErrInvalidImageFormat
		}
	} else {
		// Apply quality to existing format
		if imageType, exists := mimeTypeMap[contentType]; exists {
			options.Type = imageType
		}
	}

	// Process image
	processedImage, err := img.Process(options)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", apperrors.ErrImageProcessingFailed, err)
	}

	if len(processedImage) == 0 {
		return nil, "", apperrors.ErrImageProcessingFailed
	}

	return processedImage, outputMimeType, nil
}

// GetFileExtension returns the file extension for a given MIME type
func GetFileExtension(mimeType string) string {
	switch mimeType {
	case "image/webp":
		return "webp"
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/avif":
		return "avif"
	default:
		return "jpg"
	}
}
