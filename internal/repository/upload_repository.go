package repository

import (
	"context"

	"github.com/google/uuid"
	"go-upload/internal/domain/entity"
)

// UploadRepository defines the interface for upload data access
type UploadRepository interface {
	Create(ctx context.Context, upload *entity.Upload) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Upload, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Upload, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
