package entity

import (
	"time"

	"github.com/google/uuid"
)

type Upload struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	FilePath  string
	UserID    uuid.UUID
}
