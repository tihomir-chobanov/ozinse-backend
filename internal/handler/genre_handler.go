package handler

import (
	"net/http"
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
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} model.Genre
// @Failure 500 {object} gin.H
// @Router /api/genres [get]
func (h *GenreHandler) GetAll(c *gin.Context) {
	genres, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, genres)
}

// GetByID genre by ID
// @Summary Get genre by ID
// @Description Returns a genre by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} model.Genre
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/genres/{id} [get]
func (h *GenreHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	genre, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "genre with id " + c.Param("id") + " not found"})
		return
	}
	c.JSON(http.StatusOK, genre)
}

// Create genre
// @Summary Create genre
// @Description Creates a new genre
// @Tags categories
// @Accept json
// @Produce json
// @Param genre body model.Genre true "Genre payload"
// @Success 201 {object} model.Genre
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/genres [post]
func (h *GenreHandler) Create(c *gin.Context) {
	var genre model.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&genre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, genre)
}

// Update genre
// @Summary Update genre
// @Description Updates an existing genre by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Param genre body model.Genre true "Genre payload"
// @Success 200 {object} model.Genre
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/genres/{id} [put]
func (h *GenreHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var genre model.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	genre.ID = id
	if err := h.service.Update(&genre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, genre)
}

// Delete genre
// @Summary Delete genre
// @Description Deletes a genre by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/genres/{id} [delete]
func (h *GenreHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}
