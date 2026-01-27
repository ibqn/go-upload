package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-upload/internal/handler"
	"go-upload/internal/middleware"
	"go-upload/internal/repository/postgres"
	"go-upload/internal/router"
	"go-upload/internal/service"
	"go-upload/pkg/jwt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp(t *testing.T) *gin.Engine {
	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&postgres.UserModel{}, &postgres.UploadModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Initialize dependencies
	userRepo := postgres.NewUserRepository(db)
	uploadRepo := postgres.NewUploadRepository(db)
	jwtService := jwt.NewService("test-secret")
	storageService := service.NewStorageService("test-storage")
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

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	app := gin.New()

	// Setup routes
	router.SetupRoutes(app, authHandler, uploadHandler, fileHandler, imageHandler, authMiddleware)

	return app
}

func TestAuthIntegration_SignUpAndSignIn(t *testing.T) {
	app := setupTestApp(t)

	// Test SignUp
	signUpPayload := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	signUpBody, _ := json.Marshal(signUpPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(signUpBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var signUpResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &signUpResp)
	assert.NotEmpty(t, signUpResp["accessToken"])
	assert.Equal(t, "User created successfully", signUpResp["message"])

	// Test SignIn
	signInPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	signInBody, _ := json.Marshal(signInPayload)

	req, _ = http.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBuffer(signInBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var signInResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &signInResp)
	assert.NotEmpty(t, signInResp["accessToken"])
	assert.Equal(t, "Sign-in successful", signInResp["message"])
}

func TestAuthIntegration_SignUp_DuplicateEmail(t *testing.T) {
	app := setupTestApp(t)

	// First signup
	signUpPayload := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	signUpBody, _ := json.Marshal(signUpPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(signUpBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Second signup with same email
	req, _ = http.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(signUpBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Email already in use", resp["error"])
}

func TestAuthIntegration_GetUser_RequiresAuth(t *testing.T) {
	app := setupTestApp(t)

	// Try to get user without authentication
	req, _ := http.NewRequest(http.MethodGet, "/api/auth/user", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
