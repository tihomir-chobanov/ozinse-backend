package handler

import (
	"net/http"
	"strconv"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
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
	// We get the userID from the JWT token, which is set in the context by the authentication middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Bind the JSON body to our FavoriteRequest struct
	var req model.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service to add the favorite project for the user
	err := h.favoriteService.AddFavorite(userID.(int), req.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	// 1. We get the userID from the JWT token, which is set in the context by the authentication middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// 2. We get the project_id from the URL parameter
	projectIDParam := c.Param("project_id")
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// 3. We call the service to remove the favorite project for the user
	err = h.favoriteService.RemoveFavorite(userID.(int), projectID)
	if err != nil {
		// The repository returns an error if the project is not in favorites, so we return 404
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project removed from favorites successfully"})
}

// here we used param body for the AddFavorite because we need to send the project_id in the request body, while for RemoveFavorite we used path param because we need to specify which project to remove from favorites in the URL. These different approaches are common in RESTful API design: POST requests often use body parameters to create or link resources, while DELETE requests typically use path parameters to specify which resource to delete.