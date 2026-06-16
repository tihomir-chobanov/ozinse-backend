package handler

import (
	"net/http"
	"ozinse-backend/internal/logger"
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
// @Failure 500 {object} map[string]string
// @Router /api/categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	// Define context fields for structured logging
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/categories",
	)

	categories, err := h.service.GetAll()
	if err != nil {
		// Log persistence layer or relational database connection disruptions
		log.Error("failed to retrieve all categories from database", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log a successful multi-record fetch operation execution
	log.Info("successfully fetched all categories", "count", len(categories))
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/categories/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		// Log validation errors when URL string inputs cannot map to internal integers
		log.Warn("invalid category ID path parameter format submitted", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		// Log empty lookup parameters gracefully without throwing system errors
		log.Warn("requested category record missing from storage systems", "category_id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "category with id " + idParam + " not found"})
		return
	}

	// Log successful identification showing record property keys
	log.Info("successfully located single category record", "category_id", id, "category_name", category.Name)
	c.JSON(http.StatusOK, category)
}

// Create category
// @Summary Create category (Admin only)
// @Description Creates a new category. Requires administrator role.
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category body model.Category true "Category payload"
// @Success 201 {object} model.Category
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/categories",
	)

	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		// Log parsing syntax problems or blank request context content
		log.Error("failed to unpack structural category JSON mapping definitions", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&category); err != nil {
		// Log relational constraints or sequence assignment exceptions
		log.Error("failed to process category initialization sequence transaction", "category_name", category.Name, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log full structural persistence success indicating the runtime index id
	log.Info("new category entity successfully registered and generated", "generated_id", category.ID, "category_name", category.Name)
	c.JSON(http.StatusCreated, category)
}

// Update category
// @Summary Update category (Admin only)
// @Description Updates an existing category by ID. Requires administrator role.
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body model.Category true "Category payload"
// @Success 200 {object} model.Category
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "PUT",
		"endpoint", "/api/categories/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid category ID structure supplied for record transformation", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		log.Error("failed to decode updated category configuration schema", "category_id", id, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.ID = id
	if err := h.service.Update(&category); err != nil {
		log.Error("failed to execute requested field value overrides in database system", "category_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("category specifications successfully updated in backend records", "category_id", id, "updated_name", category.Name)
	c.JSON(http.StatusOK, category)
}

// Delete category
// @Summary Delete category (Admin only)
// @Description Deletes a category by ID. Requires administrator role.
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string "successfully deleted"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "DELETE",
		"endpoint", "/api/categories/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid index parameter schema submitted for category destruction", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		log.Error("failed to delete category record row entry from database relation", "category_id", id, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Info("category resource completely purged and drop actions confirmed", "deleted_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

// GetMovieCount handles retrieving statistics for video assets distribution across modules
// @Summary Get movie count per category (Admin only)
// @Description Returns a list of categories along with the total count of associated projects. Requires administrator role.
// @Tags categories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.CategoryMovieCount
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 500 {object} map[string]string
// @Router /api/categories/movie-count [get]
func (h *CategoryHandler) GetMovieCount(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/categories/movie-count",
	)

	results, err := h.service.GetMovieCountPerCategory()
	if err != nil {
		log.Error("failed to compile relational aggregate category stats data row values", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("successfully processed telemetry asset metrics calculations sets rows", "count", len(results))
	c.JSON(http.StatusOK, results)
}