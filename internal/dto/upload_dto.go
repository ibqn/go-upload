package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateUploadRequest represents the file upload request
// File is handled via multipart/form-data, not in JSON
type CreateUploadRequest struct {
	Folder string `form:"folder"` // Optional folder parameter
}

// UploadResponse represents a single upload record (restrictive)
type UploadResponse struct {
	ID        uuid.UUID `json:"id"`
	FilePath  string    `json:"filePath"` // Keep filePath for reference
	CreatedAt time.Time `json:"createdAt"`
}

// UploadListResponse represents a list of uploads
type UploadListResponse struct {
	Message string           `json:"message"`
	Uploads []UploadResponse `json:"uploads"`
}

// CreateUploadResponse represents the response after creating an upload
type CreateUploadResponse struct {
	Message  string    `json:"message"`
	UploadID uuid.UUID `json:"uploadId"`
}
