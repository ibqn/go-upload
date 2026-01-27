package postgres

import (
	"time"

	"github.com/google/uuid"
	"go-upload/internal/domain/entity"
	"gorm.io/gorm"
)

// UserModel represents the GORM user model for database operations
type UserModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// TableName specifies the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}

// BeforeCreate hook to generate UUID before creating
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// ToEntity converts GORM UserModel to domain entity
func (u *UserModel) ToEntity() *entity.User {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return &entity.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
	}
}

// FromEntity converts domain entity to GORM UserModel
func (u *UserModel) FromEntity(user *entity.User) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	if user.DeletedAt != nil {
		u.DeletedAt = gorm.DeletedAt{Time: *user.DeletedAt, Valid: true}
	}
	u.Username = user.Username
	u.Email = user.Email
	u.Password = user.Password
}

// UploadModel represents the GORM upload model for database operations
type UploadModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FilePath string    `gorm:"not null"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
}

// TableName specifies the table name for UploadModel
func (UploadModel) TableName() string {
	return "uploads"
}

// BeforeCreate hook to generate UUID before creating
func (u *UploadModel) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// ToEntity converts GORM UploadModel to domain entity
func (u *UploadModel) ToEntity() *entity.Upload {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return &entity.Upload{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
		FilePath:  u.FilePath,
		UserID:    u.UserID,
	}
}

// FromEntity converts domain entity to GORM UploadModel
func (u *UploadModel) FromEntity(upload *entity.Upload) {
	u.ID = upload.ID
	u.CreatedAt = upload.CreatedAt
	u.UpdatedAt = upload.UpdatedAt
	if upload.DeletedAt != nil {
		u.DeletedAt = gorm.DeletedAt{Time: *upload.DeletedAt, Valid: true}
	}
	u.FilePath = upload.FilePath
	u.UserID = upload.UserID
}
