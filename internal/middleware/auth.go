package middleware

import (
	"net/http"
	"ozinse-backend/internal/auth"
	"ozinse-backend/internal/logger"
	"strings"

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
			// Log token validation failures for security tracking
			logger.Log.Warn("client request blocked due to invalid or expired credentials token",
				"requestType", c.Request.Method,
				"endpoint", c.Request.URL.Path,
				"ip_address", c.ClientIP(),
				"error", err.Error(),
			)
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

// AdminOnly restricts access exclusively to users possessing an administrator role framework configuration.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract role_id context attributes injected by AuthMiddleware
		roleIDValue, exists := c.Get("role_id")
		userID, _ := c.Get("user_id")

		if !exists {
			// Log missing claims properties inside verified payload frameworks
			logger.Log.Warn("administrative access rejected due to missing role identity context metrics",
				"requestType", c.Request.Method,
				"endpoint",    c.Request.URL.Path,
				"ip_address",  c.ClientIP(),
			)
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Role not found"})
			c.Abort()
			return
		}

		// 2. Safely assert the interface value to an explicit integer configuration type
		roleID, ok := roleIDValue.(int)
		
		// 3. Evaluate if the account role index maps to authorization constraints schema (role_id 2 = admin)
		if !ok || roleID != 2 {
			// Log the unauthorized admin resource breach invocation attempt safely
			logger.Log.Warn("unauthorized administrative access attempt blocked successfully",
				"requestType", c.Request.Method,
				"endpoint",    c.Request.URL.Path,
				"user_id",     userID,
				"role_id",     roleIDValue, // Using interface variable safely for logging
				"ip_address",  c.ClientIP(),
			)

			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Admins only"})
			c.Abort()
			return
		}
		
		// 4. Authorization bounds confirmed, step forward into transaction layers
		c.Next()
	}
}