package utils

import (
	"fmt"
	"go-upload/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDb() {
	fmt.Println("Connecting to the database...")

	var err error

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalln("DATABASE_URL is not set")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}

	DB.AutoMigrate(&models.User{}, &models.Upload{})

	fmt.Println("Database connection established.")
}
