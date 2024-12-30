package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("my-secret-key")

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract the token string
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))

		// Parse the token
		_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Signing method check process
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorSignatureInvalid)
			}
			// Return the secret key
			return secretKey, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
			c.Abort()
			return
		}

		// Token is valid, continue processing the request
		c.Next()
	}
}
