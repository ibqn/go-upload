package middleware

import (
	"go-upload/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	accessToken, err := c.Cookie("accessToken")

	if err != nil || accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			accessToken = authHeader[7:]
		}
	}

	if accessToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "No access token provided",
		})
		return
	}

	claims, err := utils.ValidateToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid access token",
		})
		return
	}

	c.Set("userId", claims.UserId)
	c.Next()
}
