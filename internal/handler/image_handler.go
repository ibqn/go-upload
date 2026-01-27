package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "go-upload/internal/domain/errors"
	"go-upload/internal/dto"
	"go-upload/internal/service"
)

// ImageHandler handles image processing HTTP requests
type ImageHandler struct {
	imageService service.ImageService
}

// NewImageHandler creates a new image handler
func NewImageHandler(imageService service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

// GetImage handles GET /image/:id
func (h *ImageHandler) GetImage(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Parse query parameters
	params := parseImageParams(c)

	// Validate parameters
	if err := h.imageService.ValidateImageParams(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image parameters"})
		return
	}

	// Process image
	processedImage, mimeType, err := h.imageService.ProcessImage(c.Request.Context(), fileID, params)
	if err != nil {
		handleImageError(c, err)
		return
	}

	// Set response headers
	c.Header("Content-Type", mimeType)
	ext := service.GetFileExtension(mimeType)
	filename := fmt.Sprintf("image_%s.%s", fileID.String(), ext)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))

	// Send image data
	c.Data(http.StatusOK, mimeType, processedImage)
}

func parseImageParams(c *gin.Context) dto.ImageParams {
	params := dto.ImageParams{
		Quality: 80, // Default quality
	}

	// Parse width
	if w := c.Query("w"); w != "" {
		if width, err := strconv.Atoi(w); err == nil && width > 0 {
			params.Width = width
		}
	}

	// Parse quality
	if q := c.Query("q"); q != "" {
		if quality, err := strconv.Atoi(q); err == nil && quality > 0 && quality <= 100 {
			params.Quality = quality
		}
	}

	// Parse format
	params.Format = c.Query("format")

	return params
}

func handleImageError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrFileNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
	case errors.Is(err, apperrors.ErrInvalidFileType):
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is not an image"})
	case errors.Is(err, apperrors.ErrInvalidImageFormat):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
	case errors.Is(err, apperrors.ErrInvalidImageParams):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image parameters"})
	case errors.Is(err, apperrors.ErrImageProcessingFailed):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
