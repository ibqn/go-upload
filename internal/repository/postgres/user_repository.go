package postgres

import (
	"context"

	"github.com/google/uuid"
	"go-upload/internal/domain/entity"
	"go-upload/internal/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	model := &UserModel{}
	model.FromEntity(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Update entity with generated ID
	*user = *model.ToEntity()
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	model := &UserModel{}
	model.FromEntity(user)

	return r.db.WithContext(ctx).Save(model).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id).Error
}
