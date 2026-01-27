package service

import (
	"context"
	"errors"
	"mime"
	"path/filepath"

	"github.com/google/uuid"
	apperrors "go-upload/internal/domain/errors"
	"go-upload/internal/repository"
	"gorm.io/gorm"
)

// fileService handles file retrieval operations
type fileService struct {
	uploadRepo repository.UploadRepository
}

// NewFileService creates a new file service
func NewFileService(uploadRepo repository.UploadRepository) FileService {
	return &fileService{
		uploadRepo: uploadRepo,
	}
}

// GetFile retrieves a file by ID
func (s *fileService) GetFile(ctx context.Context, fileID uuid.UUID) (string, string, error) {
	upload, err := s.uploadRepo.FindByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", apperrors.ErrFileNotFound
		}
		return "", "", err
	}

	mimeType := s.DetectMimeType(upload.FilePath)
	return upload.FilePath, mimeType, nil
}

// DetectMimeType detects the MIME type of a file based on its extension
func (s *fileService) DetectMimeType(filePath string) string {
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	return mimeType
}
