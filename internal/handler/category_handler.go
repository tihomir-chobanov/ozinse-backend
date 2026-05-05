package handler

import (
	"net/http"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// GetAll categories
// @Summary List categories
// @Description Returns all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} model.Category
// @Failure 500 {object} gin.H
// @Router /api/categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetByID category by ID
// @Summary Get category by ID
// @Description Returns a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.Category
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	category, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category with id " + c.Param("id") + " not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// Create category
// @Summary Create category
// @Description Creates a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body model.Category true "Category payload"
// @Success 201 {object} model.Category
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}

// Update category
// @Summary Update category
// @Description Updates an existing category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body model.Category true "Category payload"
// @Success 200 {object} model.Category
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category.ID = id
	if err := h.service.Update(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// Delete category
// @Summary Delete category
// @Description Deletes a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
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
