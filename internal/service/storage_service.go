package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/xid"
)

// storageService handles file storage operations
type storageService struct {
	basePath string
}

// NewStorageService creates a new storage service
func NewStorageService(basePath string) StorageService {
	return &storageService{
		basePath: basePath,
	}
}

// SaveFile saves an uploaded file to the filesystem
// Returns the full path where the file was saved
func (s *storageService) SaveFile(file *multipart.FileHeader, userID uuid.UUID, folder string) (string, error) {
	// Clean folder path to prevent directory traversal
	cleanFolder := filepath.Clean(folder)
	storagePath := filepath.Join(s.basePath, userID.String(), cleanFolder)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Get original filename
	fileName := file.Filename
	fullPath := filepath.Join(storagePath, fileName)

	// Handle filename conflicts by appending unique ID
	if _, err := os.Stat(fullPath); err == nil {
		ext := filepath.Ext(fileName)
		nameWithoutExt := strings.TrimSuffix(fileName, ext)
		uniqueID := xid.New().String()
		fileName = fmt.Sprintf("%s_%s%s", nameWithoutExt, uniqueID, ext)
		fullPath = filepath.Join(storagePath, fileName)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return fullPath, nil
}

// DeleteFile deletes a file from the filesystem
func (s *storageService) DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// FileExists checks if a file exists at the given path
func (s *storageService) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// GetFilePath returns the full path for a given relative path
func (s *storageService) GetFilePath(relativePath string) string {
	return filepath.Join(s.basePath, relativePath)
}
