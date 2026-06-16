package handler

import (
	"net/http"
	"ozinse-backend/internal/logger"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"

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
// @Summary Create a new role with permissions (Admin only)
// @Description Inserts a new role name and maps specified module permissions. Requires administrator role.
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role body model.Role true "Role and Permissions Data"
// @Success 201 {object} model.Role
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string "Invalid JSON input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/roles",
	)

	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to decode incoming role allocation JSON definition parameters", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateRole(&req); err != nil {
		log.Error("rbac service components failed to register new authorization schema row", "name", req.Name, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("administrative authorization access role successfully initialized", "generated_id", req.ID, "name", req.Name)
	c.JSON(http.StatusCreated, req)
}

// GetAllRoles retrieves all configured roles with their grouped permissions.
// @Summary Get all roles (Admin only)
// @Description Returns a full list of roles, each containing its associated module permissions array. Requires administrator role.
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Role
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/roles [get]
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/roles",
	)

	roles, err := h.service.GetAllRoles()
	if err != nil {
		log.Error("failed to extract structured system roles collection metrics from matrix", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("successfully parsed all operational permission metadata layers", "count", len(roles))
	c.JSON(http.StatusOK, roles)
}

// GetRoleByID retrieves a specific role config by its ID.
// @Summary Get a role by ID (Admin only)
// @Description Returns a single role object containing its mapped module permissions array. Requires administrator role.
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} model.Role
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error or role not found"
// @Router /api/roles/{id} [get]
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	idStr := c.Param("id")
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/roles/"+idStr,
	)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid role sequence criteria parsing exception triggered", "raw_id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	role, err := h.service.GetRoleByID(id)
	if err != nil {
		log.Error("failed to evaluate isolated database row definition parameters for role entry", "role_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("successfully located single role structure assignment schemas", "role_id", id, "name", role.Name)
	c.JSON(http.StatusOK, role)
}

// UpdateRole updates an existing role's name and completely overrides its permissions.
// @Summary Update an existing role (Admin only)
// @Description Modifies role details by ID and replaces old permissions with the new payload within a transaction. Requires administrator role.
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body model.Role true "Updated Role Data"
// @Success 200 {object} model.Role
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string "Invalid ID or JSON input"
// @Failure 500 {object} map[string]string "Internal server error or role not found"
// @Router /api/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	log := logger.Log.With(
		"requestType", "PUT",
		"endpoint", "/api/roles/"+idStr,
	)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid conversion path parameter index matching criteria submitted for updating", "raw_id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse updated role configuration definitions content body", "role_id", id, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id // Ensure the struct ID matches the URL path parameter

	if err := h.service.UpdateRole(&req); err != nil {
		log.Error("rbac transactional modification override failed to synchronize changes", "role_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("role authorization matrices successfully altered and written", "role_id", id, "updated_name", req.Name)
	c.JSON(http.StatusOK, req)
}

// DeleteRole removes a role from the database.
// @Summary Delete a role (Admin only)
// @Description Deletes a role record by ID. Foreign key cascade constraints handle removing associated permissions automatically. Requires administrator role.
// @Tags roles
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]string "Role deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error or role not found"
// @Router /api/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	log := logger.Log.With(
		"requestType", "DELETE",
		"endpoint", "/api/roles/"+idStr,
	)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid syntax structure parsing error verified upon row purge invocation", "raw_id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	if err := h.service.DeleteRole(id); err != nil {
		log.Error("failed to safely eliminate target role index bounds from runtime registry maps", "role_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("role indices and associated constraint matrices dropped completely", "deleted_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "Role and its associated permissions deleted successfully"})
}