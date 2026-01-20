package models

import "gorm.io/gorm"

type Upload struct {
	gorm.Model

	FilePath string `gorm:"not null"`
	UserID   string `gorm:"not null"`
}
