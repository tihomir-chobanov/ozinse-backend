package handler

import (
	"net/http"
	"ozinse-backend/internal/auth"
	"ozinse-backend/internal/logger"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService    *service.AuthService
	jwtSecret      string
	jwtExpiryHours int
}

func NewAuthHandler(authService *service.AuthService, secret string, expiry int) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
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
// @Param login body model.LoginRequest true "Login Credentials"
// @Success 200 {object} map[string]string "Returns JWT token"
// @Failure 400 {object} map[string]string "Invalid input format"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Failure 500 {object} map[string]string "Failed to generate access token"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/auth/login",
	)

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to bind incoming login JSON input format", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Fetch the real user from the database via the Service
	user, err := h.authService.GetByEmail(req.Email)
	if err != nil || user == nil {
		// Log authentication warnings cleanly without exposing account existence details to clients
		log.Warn("authentication failed: invalid email or user does not exist", "email", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the hashed password from the DB with the provided plain password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Warn("authentication failed: password mismatch for existing account", "email", req.Email, "user_id", user.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a real JWT token using the validated user's data
	token, err := auth.GenerateToken(user.ID, user.Email, user.RoleID, h.jwtSecret, h.jwtExpiryHours)
	if err != nil {
		log.Error("security runtime failed to generate authorization token", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	log.Info("user successfully authenticated and issued session token", "user_id", user.ID, "role_id", user.RoleID)
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
// @Failure 400 {object} map[string]string "Invalid JSON or short password"
// @Failure 500 {object} map[string]string "Internal server error or email already exists"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/auth/register",
	)

	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to unpack structural registration request layout configurations", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service layer to hash password, enrich data, and save user
	user, err := h.authService.Register(&req)
	if err != nil {
		log.Error("identity service failed to complete new user account registration lifecycle", "email", req.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("fresh new user identity account successfully registered", "user_id", user.ID, "email", user.Email)
	c.JSON(http.StatusCreated, user)
}

// ForgotPassword generates a secure recovery token and returns it (simulating an email send).
// @Summary Request a password reset token
// @Description Generates a 60-minute expiration token linked to the user's email.
// @Tags auth
// @Accept json
// @Produce json
// @Param forgot body model.ForgotPasswordRequest true "User Email"
// @Success 200 {object} map[string]string "Token generated successfully"
// @Failure 400 {object} map[string]string "Invalid email format"
// @Failure 500 {object} map[string]string "User not found or DB error"
// @Router /api/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/auth/forgot-password",
	)

	var req model.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to validate password recovery initial body values", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate token and save expiration state in the database
	token, err := h.authService.ForgotPassword(&req)
	if err != nil {
		log.Warn("failed to build recovery verification sequence for submitted email data", "email", req.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("account security reset recovery sequence validation token generated", "email", req.Email)
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
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {object} map[string]string "Invalid JSON input"
// @Failure 500 {object} map[string]string "Invalid/expired token or server error"
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/auth/reset-password",
	)

	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to read input fields for executing token override updates", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply the new password if token status checks pass
	if err := h.authService.ResetPassword(&req); err != nil {
		log.Error("failed to process password parameter alterations via supplied sequence token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("user credential parameters modified and updated via valid validation token")
	c.JSON(http.StatusOK, gin.H{"message": "Your password has been reset successfully"})
}