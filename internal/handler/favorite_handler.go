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

// RemoveFavorite handles DELETE /api/users/favorites/:project_id
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