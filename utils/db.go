package utils

import (
	"fmt"
	"go-upload/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDb() {
	fmt.Println("Connecting to the database...")

	var err error

	DB, err = gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}

	DB.AutoMigrate(&models.User{})

	fmt.Println("Database connection established.")
}
