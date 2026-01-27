package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-upload/internal/service"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	authService service.AuthService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth is a middleware that validates JWT tokens
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from cookie first
		accessToken, err := c.Cookie("accessToken")

		// If not in cookie, check Authorization header
		if err != nil || accessToken == "" {
			authHeader := c.GetHeader("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				accessToken = authHeader[7:]
			}
		}

		// If no token found, return unauthorized
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No access token provided",
			})
			return
		}

		// Validate token
		userID, err := m.authService.ValidateToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid access token",
			})
			return
		}

		// Set user ID in context
		c.Set("userId", userID.String())
		c.Next()
	}
}
