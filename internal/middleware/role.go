package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// AdminOnly checks if the user has the administrator role.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract role_id that was set by AuthMiddleware
		roleIDValue, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Role not found"})
			c.Abort()
			return
		}

		// 2. Safely assert the type to integer
		roleID, ok := roleIDValue.(int)
		
		// 3. Check if the role matches the admin ID (which is 2 in your database)
		if !ok || roleID != 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Admins only"})
			c.Abort()
			return
		}
		
		// Allow the request to proceed if checks pass
		c.Next()
	}
}