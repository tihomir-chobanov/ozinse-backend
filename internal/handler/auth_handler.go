package handler

import (
	"net/http"
	"ozinse-backend/internal/auth"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)



type AuthHandler struct {
	authService *service.AuthService
	jwtSecret      string
	jwtExpiryHours int
}

func NewAuthHandler(authService *service.AuthService, secret string, expiry int) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:      secret,
		jwtExpiryHours: expiry,
	}
}

// Login handles the user authentication process.
// @Summary User login
// @Description Authenticates user credentials and returns a secure JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "Login Credentials" // <-- ПРОМЯНА: става model.LoginRequest
// @Success 200 {object} gin.H "Returns JWT token"
// @Failure 400 {object} gin.H "Invalid input format"
// @Failure 401 {object} gin.H "Invalid email or password"
// @Failure 500 {object} gin.H "Failed to generate access token"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req  model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// 1. Fetch the real user from the database via the Service
	user, err := h.authService.GetByEmail(req.Email)
	// NEW: Check BOTH if there was an error OR if the user simply wasn't found (is nil)
	if err != nil || user == nil {
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

// Register handles user signup and creates a profile with default values.
// @Summary Register a new user
// @Description Creates a new user with email and password. Other fields get default values.
// @Tags auth
// @Accept json
// @Produce json
// @Param register body model.RegisterRequest true "Registration Credentials"
// @Success 201 {object} model.User
// @Failure 400 {object} gin.H "Invalid JSON or short password"
// @Failure 500 {object} gin.H "Internal server error or email already exists"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	// Validate incoming JSON structure and character length constraints
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service layer to hash password, enrich data, and save user
	user, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}



// ForgotPassword generates a secure recovery token and returns it (simulating an email send).
// @Summary Request a password reset token
// @Description Generates a 60-minute expiration token linked to the user's email.
// @Tags auth
// @Accept json
// @Produce json
// @Param forgot body model.ForgotPasswordRequest true "User Email"
// @Success 200 {object} gin.H "Token generated successfully"
// @Failure 400 {object} gin.H "Invalid email format"
// @Failure 500 {object} gin.H "User not found or DB error"
// @Router /api/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req model.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate token and save expiration state in the database
	token, err := h.authService.ForgotPassword(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// In production, this token goes into an email link. For now, we return it to test in Postman.
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset token generated successfully. In production, this is sent via email.",
		"token":   token,
	})
}

// ResetPassword validates the recovery token and saves the new password configuration.
// @Summary Reset password using token
// @Description Verifies if the token is valid and not expired, then updates the user's password.
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body model.ResetPasswordRequest true "Token and New Password"
// @Success 200 {object} gin.H "Password updated successfully"
// @Failure 400 {object} gin.H "Invalid JSON input"
// @Failure 500 {object} gin.H "Invalid/expired token or server error"
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply the new password if token status checks pass
	if err := h.authService.ResetPassword(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your password has been reset successfully"})
}
