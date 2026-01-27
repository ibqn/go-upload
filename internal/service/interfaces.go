package service

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"go-upload/internal/dto"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	SignUp(ctx context.Context, req *dto.SignUpRequest) (*dto.AuthResponse, error)
	SignIn(ctx context.Context, req *dto.SignInRequest) (*dto.AuthResponse, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
}

// UploadService defines the interface for upload operations
type UploadService interface {
	CreateUpload(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID, folder string) (*dto.CreateUploadResponse, error)
	ListUploads(ctx context.Context, userID uuid.UUID) (*dto.UploadListResponse, error)
	GetUpload(ctx context.Context, uploadID, userID uuid.UUID) (*dto.UploadResponse, error)
	DeleteUpload(ctx context.Context, uploadID, userID uuid.UUID) error
}

// FileService defines the interface for file retrieval operations
type FileService interface {
	GetFile(ctx context.Context, fileID uuid.UUID) (string, string, error)
	DetectMimeType(filePath string) string
}

// ImageService defines the interface for image processing operations
type ImageService interface {
	ProcessImage(ctx context.Context, fileID uuid.UUID, params dto.ImageParams) ([]byte, string, error)
	ValidateImageParams(params *dto.ImageParams) error
}

// StorageService defines the interface for file storage operations
type StorageService interface {
	SaveFile(file *multipart.FileHeader, userID uuid.UUID, folder string) (string, error)
	DeleteFile(filePath string) error
	FileExists(filePath string) bool
	GetFilePath(relativePath string) string
}
