package errors

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidToken       = errors.New("invalid token")
	ErrEmailAlreadyExists = errors.New("email already exists")

	// User errors
	ErrUserNotFound = errors.New("user not found")

	// Upload errors
	ErrUploadNotFound       = errors.New("upload not found")
	ErrFileRequired         = errors.New("file is required")
	ErrFileTooLarge         = errors.New("file size exceeds maximum allowed (10MB)")
	ErrInvalidFileType      = errors.New("invalid file type")
	ErrUploadFailed         = errors.New("file upload failed")
	ErrFileNotFound         = errors.New("file not found")
	ErrUnauthorizedAccess   = errors.New("unauthorized access to resource")

	// Image errors
	ErrInvalidImageFormat = errors.New("invalid image format")
	ErrInvalidImageParams = errors.New("invalid image parameters")
	ErrImageProcessingFailed = errors.New("image processing failed")

	// Validation errors
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidUUID     = errors.New("invalid UUID format")
	ErrMissingRequired = errors.New("missing required field")

	// Storage errors
	ErrStorageFailed = errors.New("storage operation failed")
	ErrFileDeleteFailed = errors.New("file deletion failed")
)
