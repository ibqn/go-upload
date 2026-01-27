package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-upload/internal/domain/entity"
	apperrors "go-upload/internal/domain/errors"
	"go-upload/internal/dto"
	"go-upload/internal/service"
	"go-upload/pkg/jwt"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	// Simulate ID generation
	user.ID = uuid.New()
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestAuthService_SignUp_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := jwt.NewService("test-secret")
	authService := service.NewAuthService(mockRepo, jwtService)
	ctx := context.Background()

	req := &dto.SignUpRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock: Email doesn't exist
	mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)
	// Mock: User creation succeeds
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	resp, err := authService.SignUp(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, "User created successfully", resp.Message)
	assert.Equal(t, req.Username, resp.User.Username)
	assert.Equal(t, req.Email, resp.User.Email)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_SignUp_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := jwt.NewService("test-secret")
	authService := service.NewAuthService(mockRepo, jwtService)
	ctx := context.Background()

	req := &dto.SignUpRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	existingUser := &entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Username: "existing",
	}

	// Mock: Email already exists
	mockRepo.On("FindByEmail", ctx, req.Email).Return(existingUser, nil)

	resp, err := authService.SignUp(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apperrors.ErrEmailAlreadyExists))

	mockRepo.AssertExpectations(t)
}

func TestAuthService_SignIn_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := jwt.NewService("test-secret")
	authService := service.NewAuthService(mockRepo, jwtService)
	ctx := context.Background()

	// Pre-hashed password: "password123"
	hashedPassword := "$2a$10$K6BYlK00xuFHws9SrAkrqO97MZKHA4oSfSA6akBVkYZc1PpoKmkI2"

	existingUser := &entity.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Username: "testuser",
		Password: hashedPassword,
	}

	req := &dto.SignInRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock: User exists
	mockRepo.On("FindByEmail", ctx, req.Email).Return(existingUser, nil)

	resp, err := authService.SignIn(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, "Sign-in successful", resp.Message)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_SignIn_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtService := jwt.NewService("test-secret")
	authService := service.NewAuthService(mockRepo, jwtService)
	ctx := context.Background()

	req := &dto.SignInRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Mock: User not found
	mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	resp, err := authService.SignIn(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apperrors.ErrInvalidCredentials))

	mockRepo.AssertExpectations(t)
}
