package postgres

import (
	"context"

	"github.com/google/uuid"
	"go-upload/internal/domain/entity"
	"go-upload/internal/repository"
	"gorm.io/gorm"
)

type uploadRepository struct {
	db *gorm.DB
}

// NewUploadRepository creates a new instance of UploadRepository
func NewUploadRepository(db *gorm.DB) repository.UploadRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) Create(ctx context.Context, upload *entity.Upload) error {
	model := &UploadModel{}
	model.FromEntity(upload)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Update entity with generated ID
	*upload = *model.ToEntity()
	return nil
}

func (r *uploadRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Upload, error) {
	var model UploadModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *uploadRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Upload, error) {
	var models []UploadModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	uploads := make([]entity.Upload, len(models))
	for i, model := range models {
		uploads[i] = *model.ToEntity()
	}

	return uploads, nil
}

func (r *uploadRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&UploadModel{}, "id = ?", id).Error
}
