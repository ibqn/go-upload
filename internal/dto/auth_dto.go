package dto

import "github.com/google/uuid"

// SignUpRequest represents the signup request payload
type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// SignInRequest represents the signin request payload
type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the authentication response with token
type AuthResponse struct {
	Message     string       `json:"message"`
	AccessToken string       `json:"accessToken"`
	User        UserResponse `json:"user,omitempty"`
}

// UserResponse represents the user information (restrictive - no password, no timestamps)
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}
