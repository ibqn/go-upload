package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "go-upload/internal/domain/errors"
	"go-upload/internal/dto"
	"go-upload/internal/service"
)

// UploadHandler handles upload HTTP requests
type UploadHandler struct {
	uploadService service.UploadService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(uploadService service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// CreateUpload handles POST /api/upload/
func (h *UploadHandler) CreateUpload(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Get optional folder parameter
	folder := c.DefaultPostForm("folder", "")

	// Create upload
	resp, err := h.uploadService.CreateUpload(c.Request.Context(), file, userID, folder)
	if err != nil {
		handleUploadError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListUploads handles GET /api/upload/
func (h *UploadHandler) ListUploads(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.uploadService.ListUploads(c.Request.Context(), userID)
	if err != nil {
		handleUploadError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUpload handles GET /api/upload/:id
func (h *UploadHandler) GetUpload(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uploadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload ID"})
		return
	}

	resp, err := h.uploadService.GetUpload(c.Request.Context(), uploadID, userID)
	if err != nil {
		handleUploadError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload retrieved successfully",
		"upload":  resp,
	})
}

// DeleteUpload handles DELETE /api/upload/:id
func (h *UploadHandler) DeleteUpload(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uploadID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload ID"})
		return
	}

	err = h.uploadService.DeleteUpload(c.Request.Context(), uploadID, userID)
	if err != nil {
		handleUploadError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Upload deleted successfully"})
}

func handleUploadError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrFileTooLarge):
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
	case errors.Is(err, apperrors.ErrFileRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
	case errors.Is(err, apperrors.ErrUploadNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
	case errors.Is(err, apperrors.ErrUnauthorizedAccess):
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
	case errors.Is(err, apperrors.ErrUploadFailed):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
	case errors.Is(err, apperrors.ErrFileDeleteFailed):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
