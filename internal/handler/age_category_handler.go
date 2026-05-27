package handler

import (
	"net/http"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AgeCategoryHandler struct {
	service *service.AgeCategoryService
}

func NewAgeCategoryHandler(service *service.AgeCategoryService) *AgeCategoryHandler {
	return &AgeCategoryHandler{service: service}
}

// GetAll age_categories
// @Summary List age_categories
// @Description Returns all age_categories
// @Tags age_categories
// @Accept json
// @Produce json
// @Success 200 {array} model.Age_Category
// @Failure 500 {object} gin.H
// @Router /api/categories [get]
func (h *AgeCategoryHandler) GetAll(c *gin.Context) {
	age_categories, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, age_categories)
}

// GetByID age_category by ID
// @Summary Get age_category by ID
// @Description Returns an age_category by its ID
// @Tags age_categories
// @Accept json
// @Produce json
// @Param id path int true "Age_Category ID"
// @Success 200 {object} model.Age_Category
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/categories/{id} [get]
func (h *AgeCategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	age_category, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "age_category with id " + c.Param("id") + " not found"})
		return
	}
	c.JSON(http.StatusOK, age_category)
}

// Create age_category
// @Summary Create age_category
// @Description Creates a new age_category
// @Tags age_categories
// @Accept json
// @Produce json
// @Param category body model.Age_Category true "Age_Category payload"
// @Success 201 {object} model.Age_Category
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/categories [post]
func (h *AgeCategoryHandler) Create(c *gin.Context) {
	var age_category model.Age_Category
	if err := c.ShouldBindJSON(&age_category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&age_category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, age_category)
}

// Update age_category
// @Summary Update age_category
// @Description Updates an existing age_category by ID
// @Tags age_categories
// @Accept json
// @Produce json
// @Param id path int true "Age_Category ID"
// @Param category body model.Age_Category true "Age_Category payload"
// @Success 200 {object} model.Age_Category
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/categories/{id} [put]
func (h *AgeCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var age_category model.Age_Category
	if err := c.ShouldBindJSON(&age_category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	age_category.ID = id
	if err := h.service.Update(&age_category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, age_category)
}

// Delete age_category
// @Summary Delete age_category
// @Description Deletes an age_category by ID
// @Tags age_categories
// @Accept json
// @Produce json
// @Param id path int true "Age_Category ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/categories/{id} [delete]
func (h *AgeCategoryHandler) Delete(c *gin.Context) {
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
