package utils

import (
	"go-upload/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User) (string, error) {
	claims := &Claims{
		UserId: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
