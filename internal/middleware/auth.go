package middleware

import (
	"net/http"
	"strings"
	"ozinse-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepts HTTP requests to verify the presence and validity of a JWT token.
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract the Authorization header from the request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 2. Ensure the header follows the "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be formatted as 'Bearer <token>'"})
			c.Abort()
			return
		}

		// 3. Validate the extracted token using the auth package
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 4. Store the extracted user data in the Gin context
		// This allows subsequent handlers to know exactly who is making the request
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID) 

		// 5. Allow the request to proceed to the next handler
		c.Next()
	}
}