package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-upload/internal/domain/entity"
	"go-upload/internal/repository/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&postgres.UserModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Create user
	user := &entity.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Find by email
	foundUser, err := repo.FindByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Create user
	user := &entity.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Find by ID
	foundUser, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Try to find non-existent user
	foundUser, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Nil(t, foundUser)
}
