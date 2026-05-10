package handler

import (
	"net/http"
	"ozinse-backend/internal/auth"
	"ozinse-backend/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	userService    *service.UserService
	jwtSecret      string
	jwtExpiryHours int
}

func NewAuthHandler(us *service.UserService, secret string, expiry int) *AuthHandler {
	return &AuthHandler{
		userService:    us,
		jwtSecret:      secret,
		jwtExpiryHours: expiry,
	}
}

// Login handles the user authentication process.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// 1. Fetch the real user from the database via the Service
	user, err := h.userService.GetByEmail(req.Email)
	if err != nil {
		// Generic message for security: don't reveal if the email exists or not
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 2. Compare the hashed password from the DB with the provided plain password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 3. Generate a real JWT token using the validated user's data
	token, err := auth.GenerateToken(user.ID, user.Email, user.RoleID, h.jwtSecret, h.jwtExpiryHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}