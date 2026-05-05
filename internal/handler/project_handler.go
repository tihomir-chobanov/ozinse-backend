package handler

import (
	"net/http"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateProjectDTO struct {
	model.Project
	GenreIDs       []int `json:"genre_ids"`
	AgeCategoryIDs []int `json:"age_category_ids"`
	CategoryIDs    []int `json:"category_ids"`
}

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(service *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// GetAll projects
// @Summary List projects
// @Description Returns all projects
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} model.Project
// @Failure 500 {object} gin.H
// @Router /api/projects [get]
func (h *ProjectHandler) GetAll(c *gin.Context) {
	projects, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

// GetByID project by ID
// @Summary Get project by ID
// @Description Returns a project by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} model.Project
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	project, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project with id " + c.Param("id") + " not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

// Create project
// @Summary Create project
// @Description Creates a new project
// @Tags categories
// @Accept json
// @Produce json
// @Param project body CreateProjectDTO true "Project payload"
// @Success 201 {object} model.Project
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var req CreateProjectDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Create(&req.Project, req.GenreIDs, req.AgeCategoryIDs, req.CategoryIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, req.Project)
}

// Update project
// @Summary Update project
// @Description Updates an existing project by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body model.Project true "Project payload"
// @Success 200 {object} model.Project
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project.ID = id
	if err := h.service.Update(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

// Delete project
// @Summary Delete project
// @Description Deletes a project by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
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