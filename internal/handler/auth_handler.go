package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "go-upload/internal/domain/errors"
	"go-upload/internal/dto"
	"go-upload/internal/service"
	"gorm.io/gorm"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SignUp handles POST /api/auth/signup
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := h.authService.SignUp(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	// Set cookie
	c.SetCookie("accessToken", resp.AccessToken, 60*60*24*7, "/", "localhost", false, false)
	c.JSON(http.StatusOK, resp)
}

// SignIn handles POST /api/auth/signin
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := h.authService.SignIn(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	// Set cookie
	c.SetCookie("accessToken", resp.AccessToken, 60*60*24*7, "/", "localhost", false, false)
	c.JSON(http.StatusOK, resp)
}

// SignOut handles POST /api/auth/signout
func (h *AuthHandler) SignOut(c *gin.Context) {
	// Clear cookie
	c.SetCookie("accessToken", "", -1, "/", "localhost", false, false)
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Sign-out successful"})
}

// GetUser handles GET /api/auth/user
func (h *AuthHandler) GetUser(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// getUserIDFromContext extracts user ID from Gin context
func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, errors.New("user not authenticated")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID type")
	}

	return uuid.Parse(userIDStr)
}

// handleError converts service errors to HTTP responses
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case errors.Is(err, apperrors.ErrEmailAlreadyExists):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
	case errors.Is(err, apperrors.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case errors.Is(err, apperrors.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
