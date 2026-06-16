package handler

import (
	"net/http"
	"ozinse-backend/internal/logger"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GenreHandler struct {
	service *service.GenreService
}

func NewGenreHandler(service *service.GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

// GetAll genres
// @Summary List genres
// @Description Returns all genres
// @Tags genres
// @Accept json
// @Produce json
// @Success 200 {array} model.Genre
// @Failure 500 {object} map[string]string
// @Router /api/genres [get]
func (h *GenreHandler) GetAll(c *gin.Context) {
	// Define context fields for structured logging
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/genres",
	)

	genres, err := h.service.GetAll()
	if err != nil {
		// Log internal database operation failures
		log.Error("failed to retrieve all genres from database", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log a successful multi-record retrieval event
	log.Info("successfully retrieved all genres", "count", len(genres))
	c.JSON(http.StatusOK, genres)
}

// GetByID genre by ID
// @Summary Get genre by ID
// @Description Returns a genre by its ID
// @Tags genres
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} model.Genre
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/genres/{id} [get]
func (h *GenreHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/genres/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		// Log validation failure when URL path variables are malformed
		log.Warn("invalid genre ID path parameter format provided", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	genre, err := h.service.GetByID(id)
	if err != nil {
		// Log missing resource entities without crashing or causing an internal error
		log.Warn("requested genre record not found in persistence layer", "genre_id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "genre with id " + idParam + " not found"})
		return
	}

	// Log successful retrieval containing specific object context fields
	log.Info("successfully fetched targeted genre record", "genre_id", id, "genre_name", genre.Name)
	c.JSON(http.StatusOK, genre)
}

// Create genre
// @Summary Create genre (Admin only)
// @Description Creates a new genre. Requires administrator role.
// @Tags genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param genre body model.Genre true "Genre payload"
// @Success 201 {object} model.Genre
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/genres [post]
func (h *GenreHandler) Create(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/genres",
	)

	var genre model.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		// Log broken JSON schemas, EOF parsing errors, or mismatched fields
		log.Error("failed to bind incoming genre JSON payload structure", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&genre); err != nil {
		// Log sequence collision problems or transactional record crashes
		log.Error("service operation failed to persist new genre record", "genre_name", genre.Name, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log absolute persistence success highlighting the generated serial ID
	log.Info("new genre entity successfully registered and persisted", "generated_id", genre.ID, "genre_name", genre.Name)
	c.JSON(http.StatusCreated, genre)
}

// Update genre
// @Summary Update genre (Admin only)
// @Description Updates an existing genre by ID. Requires administrator role.
// @Tags genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Param genre body model.Genre true "Genre payload"
// @Success 200 {object} model.Genre
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/genres/{id} [put]
func (h *GenreHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "PUT",
		"endpoint", "/api/genres/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid genre ID path parameter format provided for modification", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var genre model.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		log.Error("failed to decode updated genre modifications JSON content", "genre_id", id, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre.ID = id
	if err := h.service.Update(&genre); err != nil {
		log.Error("failed to complete targeted database updates for genre entity", "genre_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("genre definitions updated successfully within backend persistence", "genre_id", id, "updated_name", genre.Name)
	c.JSON(http.StatusOK, genre)
}

// Delete genre
// @Summary Delete genre (Admin only)
// @Description Deletes a genre by ID. Requires administrator role.
// @Tags genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} map[string]string "successfully deleted"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/genres/{id} [delete]
func (h *GenreHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "DELETE",
		"endpoint", "/api/genres/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid genre ID path parameter structure submitted for destruction", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		log.Error("failed to drop genre entry from database relation systems", "genre_id", id, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Info("genre resource completely removed and purged from backend records", "deleted_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}
