package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "go-upload/internal/domain/errors"
	"go-upload/internal/service"
)

// FileHandler handles file retrieval HTTP requests
type FileHandler struct {
	fileService service.FileService
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// GetFile handles GET /file/:id
func (h *FileHandler) GetFile(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	filePath, mimeType, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		handleFileError(c, err)
		return
	}

	c.Header("Content-Type", mimeType)
	c.File(filePath)
}

func handleFileError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrFileNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
