package handler

import (
	"net/http"
	"strconv"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
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
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
// @Failure 400 {object} map[string]string "Invalid user ID format"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	// 1. Extract the "id" parameter from the URL path (e.g., /api/users/3)
	idParam := c.Param("id")

	// 2. Convert the string parameter to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 3. Fetch the user profile from the Service layer
	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. If the user doesn't exist, return a 404 Not Found error
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 5. Sanitization: Hide the hashed password before sending the JSON response
	user.Password = ""

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
// @Failure 400 {object} map[string]string "Invalid payload or ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	// 1. Extract the "id" parameter from the URL path (e.g., /api/users/3)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id 

	err = h.service.Update(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete removes a user record from the system
// @Summary Delete a user (Admin only)
// @Description Completely purges a user record by ID from the database. Requires administrator role.
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid user ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}