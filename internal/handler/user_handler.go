package handler

import (
	"net/http"
	"ozinse-backend/internal/logger"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAll retrieves all users from the system
// @Summary List all users (Admin only)
// @Description Returns a full list of registered users. Requires administrator role.
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.User
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	// Define structured context fields for tracking the request
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/users",
	)

	users, err := h.service.GetAll()
	if err != nil {
		// Log persistent internal server issues
		log.Error("failed to retrieve all user identities from database", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log a successful bulk read operation
	log.Info("successfully fetched all registered user accounts", "count", len(users))
	c.JSON(http.StatusOK, users)
}

// GetByID retrieves a specific user profile by ID
// @Summary Get user profile by ID
// @Description Returns full metadata for a single user profile. Password hashes are automatically sanitized.
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 400 {object} map[string]string "Invalid user ID format"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/users/"+idParam,
	)

	// Convert the string parameter to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// Log validation errors when URL string inputs cannot map to internal integers
		log.Warn("invalid user identity ID path parameter format provided", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Fetch the user profile from the Service layer
	user, err := h.service.GetByID(id)
	if err != nil {
		log.Error("failed to evaluate storage records row query criteria for user profile", "user_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the user entity actually exists
	if user == nil {
		// Log missing entity lookups safely as a system warning
		log.Warn("requested user account identity record missing from persistence layer", "user_id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Sanitization: Hide the hashed password before sending the JSON response
	user.Password = ""

	// Log successful identification after ensuring object pointer is safe (not nil)
	log.Info("successfully located single sanitized user account profile", "user_id", id, "email", user.Email)
	c.JSON(http.StatusOK, user)
}

// Update modifies an existing user profile
// @Summary Update user profile
// @Description Updates user fields based on the provided payload and profile ID
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.User true "Updated User Data"
// @Success 200 {object} model.User
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 400 {object} map[string]string "Invalid payload or ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "PUT",
		"endpoint", "/api/users/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid user identity conversion target ID provided for modifications", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("failed to decode profile modifications payload data parameters struct structures", "user_id", id, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id

	err = h.service.Update(&user)
	if err != nil {
		log.Error("failed to sync profile alteration properties to backend storage records", "user_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("user account definitions and parameter values altered successfully", "user_id", id)
	c.JSON(http.StatusOK, user)
}

// Delete removes a user record from the system
// @Summary Delete a user (Admin only)
// @Description Completely purges a user record by ID from the database. Requires administrator role.
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string "Invalid user ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "DELETE",
		"endpoint", "/api/users/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid conversion formatting structure triggered during profile drop invocation", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		log.Error("failed to execute identity account termination sequence lifecycle from registry maps", "user_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("user profile permanent destruction and erasure sequence successfully confirmed", "deleted_user_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// CountUsers returns the total number of users in the system
// @Summary Count total users (Admin only)
// @Description Returns the total number of registered users. Requires administrator role.
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]int "Total user count"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users-count [get]
func (h *UserHandler) CountUsers(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/users/count",
	)
	count, err := h.service.CountUsers()
	if err != nil {
		log.Error("failed to count users", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}