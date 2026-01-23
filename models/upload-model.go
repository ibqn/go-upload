package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Upload struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FilePath string    `gorm:"not null"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
}
