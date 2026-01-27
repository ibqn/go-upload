package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go-upload/config"
	"go-upload/internal/handler"
	"go-upload/internal/middleware"
	"go-upload/internal/repository/postgres"
	"go-upload/internal/router"
	"go-upload/internal/service"
	"go-upload/pkg/jwt"
	"gorm.io/gorm"
	pgdriver "gorm.io/driver/postgres"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db := setupDatabase(cfg.DatabaseURL)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	uploadRepo := postgres.NewUploadRepository(db)

	// Initialize services
	jwtService := jwt.NewService(cfg.JWTSecret)
	storageService := service.NewStorageService(cfg.StoragePath)
	authService := service.NewAuthService(userRepo, jwtService)
	uploadService := service.NewUploadService(uploadRepo, storageService)
	fileService := service.NewFileService(uploadRepo)
	imageService := service.NewImageService(uploadRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	fileHandler := handler.NewFileHandler(fileService)
	imageHandler := handler.NewImageHandler(imageService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup Gin engine
	app := gin.Default()

	// Setup routes
	router.SetupRoutes(app, authHandler, uploadHandler, fileHandler, imageHandler, authMiddleware)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := app.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupDatabase(databaseURL string) *gorm.DB {
	log.Println("Connecting to the database...")

	db, err := gorm.Open(pgdriver.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&postgres.UserModel{}, &postgres.UploadModel{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connection established.")
	return db
}
