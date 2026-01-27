package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/google/uuid"
	"go-upload/internal/domain/entity"
	apperrors "go-upload/internal/domain/errors"
	"go-upload/internal/dto"
	"go-upload/internal/repository"
	"gorm.io/gorm"
)

const maxFileSize = 10 * 1024 * 1024 // 10MB

// uploadService handles file upload business logic
type uploadService struct {
	uploadRepo     repository.UploadRepository
	storageService StorageService
}

// NewUploadService creates a new upload service
func NewUploadService(uploadRepo repository.UploadRepository, storageService StorageService) UploadService {
	return &uploadService{
		uploadRepo:     uploadRepo,
		storageService: storageService,
	}
}

// CreateUpload handles file upload
func (s *uploadService) CreateUpload(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID, folder string) (*dto.CreateUploadResponse, error) {
	// Validate file size
	if file.Size > maxFileSize {
		return nil, apperrors.ErrFileTooLarge
	}

	// Save file to storage
	filePath, err := s.storageService.SaveFile(file, userID, folder)
	if err != nil {
		return nil, apperrors.ErrUploadFailed
	}

	// Create upload record
	upload := &entity.Upload{
		UserID:   userID,
		FilePath: filePath,
	}

	if err := s.uploadRepo.Create(ctx, upload); err != nil {
		// Try to delete the file if database operation fails
		_ = s.storageService.DeleteFile(filePath)
		return nil, err
	}

	return &dto.CreateUploadResponse{
		Message:  "File uploaded successfully",
		UploadID: upload.ID,
	}, nil
}

// ListUploads returns all uploads for a user
func (s *uploadService) ListUploads(ctx context.Context, userID uuid.UUID) (*dto.UploadListResponse, error) {
	uploads, err := s.uploadRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	uploadResponses := make([]dto.UploadResponse, len(uploads))
	for i, upload := range uploads {
		uploadResponses[i] = dto.UploadResponse{
			ID:        upload.ID,
			FilePath:  upload.FilePath,
			CreatedAt: upload.CreatedAt,
		}
	}

	return &dto.UploadListResponse{
		Message: "Uploads retrieved successfully",
		Uploads: uploadResponses,
	}, nil
}

// GetUpload retrieves a single upload by ID
func (s *uploadService) GetUpload(ctx context.Context, uploadID, userID uuid.UUID) (*dto.UploadResponse, error) {
	upload, err := s.uploadRepo.FindByID(ctx, uploadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUploadNotFound
		}
		return nil, err
	}

	// Check authorization
	if upload.UserID != userID {
		return nil, apperrors.ErrUnauthorizedAccess
	}

	return &dto.UploadResponse{
		ID:        upload.ID,
		FilePath:  upload.FilePath,
		CreatedAt: upload.CreatedAt,
	}, nil
}

// DeleteUpload deletes an upload and its associated file
func (s *uploadService) DeleteUpload(ctx context.Context, uploadID, userID uuid.UUID) error {
	// Find upload
	upload, err := s.uploadRepo.FindByID(ctx, uploadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUploadNotFound
		}
		return err
	}

	// Check authorization
	if upload.UserID != userID {
		return apperrors.ErrUnauthorizedAccess
	}

	// Delete file from storage
	if err := s.storageService.DeleteFile(upload.FilePath); err != nil {
		return apperrors.ErrFileDeleteFailed
	}

	// Delete database record
	if err := s.uploadRepo.Delete(ctx, uploadID); err != nil {
		return err
	}

	return nil
}
