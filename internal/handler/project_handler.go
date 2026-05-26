package handler

import (
	"net/http"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"
	"math"
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

// GetAll projects with pagination
// @Summary List projects with pagination
// @Description Returns a paginated list of projects along with metadata
// @Tags projects
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} gin.H "Includes 'data' array and 'pagination' object"
// @Failure 500 {object} gin.H
// @Router /api/projects [get]
func (h *ProjectHandler) GetAll(c *gin.Context) {
	// 1. Get query parameters from the URL, providing string defaults
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "10")

    // 2. Convert the string parameters to integers
    page, _ := strconv.Atoi(pageStr)
    limit, _ := strconv.Atoi(limitStr)

    // 3. Call the service layer with the pagination parameters
    projects, totalCount, err := h.service.GetAll(page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 4. Calculate the total number of pages using ceiling division
    // Example: 25 total projects / 10 per page = 2.5 -> Ceil makes it 3 pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

    // 5. Return a structured JSON response containing the data and metadata
    c.JSON(http.StatusOK, gin.H{
        "data": projects,
        "pagination": gin.H{
            "current_page":   page,
            "per_page":       limit,
            "total_projects": totalCount,
            "total_pages":    totalPages,
        },
    })
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