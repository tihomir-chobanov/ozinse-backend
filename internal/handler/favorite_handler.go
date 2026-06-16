package handler

import (
	"net/http"
	"ozinse-backend/internal/logger"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteService *service.FavoriteService
}

func NewFavoriteHandler(favoriteService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteService: favoriteService}
}

// AddFavorite adds a project to the user's favorites list
// @Summary Add project to favorites
// @Description Links a project to the authenticated user's profile as a favorite item
// @Tags favorites
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param favorite body model.FavoriteRequest true "Favorite Payload (project_id)"
// @Success 200 {object} map[string]string "Project added to favorites successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/favorites [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/users/favorites",
	)

	// Fetch the userID from the JWT token context populated by AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists {
		log.Warn("unauthorized attempt to modify favorites without valid context definitions")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Enrich our structural logger context with the verified user identifier
	log = log.With("userID", userID)

	// Bind the JSON body to our FavoriteRequest struct
	var req model.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse favorite schema input parameter configurations", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service to add the favorite project for the user
	err := h.favoriteService.AddFavorite(userID.(int), req.ProjectID)
	if err != nil {
		log.Error("persistence layer execution failed to append favorite association", "project_id", req.ProjectID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("project successfully linked and associated to user favorites", "project_id", req.ProjectID)
	c.JSON(http.StatusOK, gin.H{"message": "Project added to favorites successfully"})
}

// RemoveFavorite removes a project from the user's favorites list
// @Summary Remove project from favorites
// @Description Unlinks a project from the authenticated user's profile by project ID
// @Tags favorites
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Success 200 {object} map[string]string "Project removed from favorites successfully"
// @Failure 400 {object} map[string]string "Invalid project ID format"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 404 {object} map[string]string "Project not found in favorites"
// @Router /api/users/favorites/{project_id} [delete]
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	projectIDParam := c.Param("project_id")
	log := logger.Log.With(
		"requestType", "DELETE",
		"endpoint", "/api/users/favorites/"+projectIDParam,
	)

	// Fetch the userID from the JWT token context populated by AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists {
		log.Warn("unauthorized attempt to discard favorites item without validation parameters")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	log = log.With("userID", userID)

	// Parse the project_id from the URL path parameter
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		log.Warn("invalid target key structure provided for item removal parameters", "raw_project_id", projectIDParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Call the service to remove the favorite project for the user
	err = h.favoriteService.RemoveFavorite(userID.(int), projectID)
	if err != nil {
		// The repository returns an error if the project is not in favorites, return 404 warning log
		log.Warn("requested target resource mapping absence detected for elimination", "project_id", projectID, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Info("target association dropped and purged from user favorites records", "project_id", projectID)
	c.JSON(http.StatusOK, gin.H{"message": "Project removed from favorites successfully"})
}

// here we used param body for the AddFavorite because we need to send the project_id in the request body, while for RemoveFavorite we used path param because we need to specify which project to remove from favorites in the URL. These different approaches are common in RESTful API design: POST requests often use body parameters to create or link resources, while DELETE requests typically use path parameters to specify which resource to delete.