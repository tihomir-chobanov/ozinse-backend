package handler

import (
	"net/http"
	"strconv"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// RoleHandler handles HTTP requests for role and permission management.
type RoleHandler struct {
	service *service.RoleService
}

// NewRoleHandler creates a new RoleHandler instance.
func NewRoleHandler(service *service.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// CreateRole handles the creation of a new role along with its permissions.
// @Summary Create a new role with permissions
// @Description Inserts a new role name and maps specified module permissions
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role body model.Role true "Role and Permissions Data"
// @Success 201 {object} model.Role
// @Failure 400 {object} map[string]string "Invalid JSON input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req model.Role

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateRole(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// GetAllRoles retrieves all configured roles with their grouped permissions.
// @Summary Get all roles
// @Description Returns a full list of roles, each containing its associated module permissions array
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Role
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/roles [get]
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRoleByID retrieves a specific role config by its ID.
// @Summary Get a role by ID
// @Description Returns a single role object containing its mapped module permissions array
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} model.Role
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error or role not found"
// @Router /api/roles/{id} [get]
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	role, err := h.service.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// UpdateRole updates an existing role's name and completely overrides its permissions.
// @Summary Update an existing role
// @Description Modifies role details by ID and replaces old permissions with the new payload within a transaction
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body model.Role true "Updated Role Data"
// @Success 200 {object} model.Role
// @Failure 400 {object} map[string]string "Invalid ID or JSON input"
// @Failure 500 {object} map[string]string "Internal server error or role not found"
// @Router /api/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id // Ensure the struct ID matches the URL path parameter

	if err := h.service.UpdateRole(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// DeleteRole removes a role from the database.
// @Summary Delete a role
// @Description Deletes a role record by ID. Foreign key cascade constraints handle removing associated permissions automatically
// @Tags roles
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]string "Role deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error or role not found"
// @Router /api/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	if err := h.service.DeleteRole(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role and its associated permissions deleted successfully"})
}