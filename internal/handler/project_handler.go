package handler

import (
	"math"
	"net/http"
	"ozinse-backend/internal/logger"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
// @Param search query string false "Search projects by title"
// @Success 200 {object} object "Includes 'data' array and 'pagination' object"
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/projects [get]
func (h *ProjectHandler) GetAll(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.Query("search")

	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/projects",
		"queryPage", pageStr,
		"queryLimit", limitStr,
		"querySearch", search,
	)

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	projects, totalCount, err := h.service.GetAll(page, limit, search)
	if err != nil {
		log.Error("failed to execute paginated projects query lookup selection", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	log.Info("successfully compiled query results chunk dataset matched rows", "totalCount", totalCount, "totalPages", totalPages)
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
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} model.Project
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Project not found"
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/projects/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid project reference sequence conversion format submitted", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	project, err := h.service.GetByID(id)
	if err != nil {
		log.Warn("requested media aggregate data record absence event captured", "project_id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "project with id " + idParam + " not found"})
		return
	}

	log.Info("successfully unpacked complete nested structural aggregate metrics matching project index", "project_id", id, "title", project.Title)
	c.JSON(http.StatusOK, project)
}

// Create project
// @Summary Create project (Admin only)
// @Description Creates a new project. Requires administrator role.
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param project body model.CreateProjectDTO true "Project payload"
// @Success 201 {object} model.Project
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/projects",
	)

	var req model.CreateProjectDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to structural parse incoming complex DTO payload configuration rules", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Create(&req.Project, req.GenreIDs, req.AgeCategoryIDs, req.CategoryIDs)
	if err != nil {
		log.Error("orchestration logic service execution failure binding relational entities to transaction context", "title", req.Project.Title, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("complex video asset entities and junction mappings committed successfully", "generated_id", req.Project.ID, "title", req.Project.Title)
	c.JSON(http.StatusCreated, req.Project)
}

// Update project
// @Summary Update project (Admin only)
// @Description Updates an existing project by ID. Requires administrator role.
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body model.CreateProjectDTO true "Project payload"
// @Success 200 {object} model.Project
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "PUT",
		"endpoint", "/api/projects/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid transformation identification format mapping criteria provided", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		log.Error("failed to decode project data overrides parameters struct mapping structures", "project_id", id, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = id
	if err := h.service.Update(&project); err != nil {
		log.Error("failed to persist modifications sequence definitions on physical relational records row data", "project_id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("project record adjustments written smoothly to storage subsystem variables", "project_id", id, "updated_title", project.Title)
	c.JSON(http.StatusOK, project)
}

// Delete project
// @Summary Delete project (Admin only)
// @Description Deletes a project by ID. Requires administrator role.
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]string "successfully deleted"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "DELETE",
		"endpoint", "/api/projects/"+idParam,
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid format structure parsing error triggered upon record removal request", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		log.Error("failed to clear targeting item indices parameters within cascade dependencies operations", "project_id", id, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Info("project definitions dropped and cascade constraints metadata rows detached perfectly", "deleted_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

// GetTrending projects
// @Summary List trending projects
// @Description Returns a list of high-traffic or featured trending projects
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} model.Project
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/projects/trending [get]
func (h *ProjectHandler) GetTrending(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/projects/trending",
	)

	projects, err := h.service.GetTrending()
	if err != nil {
		log.Error("failed to extract global featured high traffic visibility datasets rows info", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("successfully processed telemetry filter results dataset fields metrics rows", "count", len(projects))
	c.JSON(http.StatusOK, projects)
}

// GetSimilar projects
// @Summary Get similar projects
// @Description Returns projects that are related to the given project ID
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} model.Project
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 400 {object} map[string]string "Invalid project ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/projects/{id}/similar [get]
func (h *ProjectHandler) GetSimilar(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/projects/"+idParam+"/similar",
	)

	projectID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid payload context target conversion processing argument format mismatch", "raw_project_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	projects, err := h.service.GetSimilar(projectID)
	if err != nil {
		log.Error("failed to run logical genre grouping lookup functions based on parameter attributes", "project_id", projectID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("successfully filtered database related items cluster categories definitions array", "source_project_id", projectID, "matched_count", len(projects))
	c.JSON(http.StatusOK, projects)
}

// CreateMainPageEntry adds a project to the main page
// @Summary Add project to main page (Admin only)
// @Description Shortcuts a project onto the featured landing main page area. Requires administrator role.
// @Tags main-page
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param entry body model.CreateMainPageEntryRequest true "Main page entry payload"
// @Success 200 {object} map[string]string "Main page entry created successfully"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/main-page-entries [post]
func (h *ProjectHandler) CreateMainPageEntry(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "POST",
		"endpoint", "/api/main-page-entries",
	)

	var req model.CreateMainPageEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to correctly unpack dashboard mapping configuration formats schemas", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateMainPageEntry(req.ProjectID, req.Position, req.IconURL); err != nil {
		log.Error("failed to mount shortcut index position layout configurations inside tables systems", "project_id", req.ProjectID, "position", req.Position, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("dashboard landing shortcut position assignment transaction confirmed completed", "project_id", req.ProjectID, "assigned_position", req.Position)
	c.JSON(http.StatusOK, gin.H{"message": "Main page entry created successfully"})
}

// GetMainPageEntries fetches all main page entries
// @Summary List main page entries (Admin only)
// @Description Returns all projects assigned to the featured main page layout. Requires administrator role.
// @Tags main-page
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} model.MainPageProject
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 500 {object} map[string]string
// @Router /api/main-page-entries [get]
func (h *ProjectHandler) GetMainPageEntries(c *gin.Context) {
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/main-page-entries",
	)

	entries, err := h.service.GetMainPageEntries()
	if err != nil {
		log.Error("failed to gather active configuration instances definitions data rows from grid maps", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("successfully aggregated total active dashboard feature shortcut slots metrics rows", "count", len(entries))
	c.JSON(http.StatusOK, entries)
}

// GetByIDForMainPage fetches a single main page entry by its configuration ID
// @Summary Get main page entry by ID (Admin only)
// @Description Returns details of a main page featured layout configuration by its database ID. Requires administrator role.
// @Tags main-page
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Main Page Entry ID"
// @Success 200 {object} model.MainPageProject
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/main-page-entries/{id} [get]
func (h *ProjectHandler) GetByIDForMainPage(c *gin.Context) {
	idParam := c.Param("id")
	log := logger.Log.With(
		"requestType", "GET",
		"endpoint", "/api/main-page-entries/"+idParam,
	)

	mainPageEntryID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid index lookup format mapping specifications structure", "raw_id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid main page entry ID"})
		return
	}

	entry, err := h.service.GetMainPageEntryById(mainPageEntryID)
	if err != nil {
		log.Error("failed to evaluate landing grid parameters criteria during runtime search functions", "entry_id", mainPageEntryID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 1. First we check if we have a record at all, because if entry is nil, it means the ID was valid but there is no such entry in the database. This is a different case than an actual error during the query execution.
	if entry == nil {
		log.Warn("targeted dashboard feature shortcut row entry does not exist", "entry_id", mainPageEntryID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Main page entry not found"})
		return
	}

	// 2. Now we can safely log, because entry is not nil
	log.Info("successfully found isolated main page grid layout metrics schema settings row", "entry_id", mainPageEntryID)
	c.JSON(http.StatusOK, entry)
}

// DeleteMainPageEntry removes a project from the main page layout
// @Summary Delete main page entry (Admin only)
// @Description Removes a featured shortcut entry from the main page by its config ID. Requires administrator role.
// @Tags main-page
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Main Page Entry ID"
// @Success 200 {object} map[string]string "Main page entry deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/main-page-entries/{id} [delete]
func (h *ProjectHandler) DeleteMainPageEntry(c *gin.Context) {
	mainPageIDStr := c.Param("id")
	log := logger.Log.With(
		"requestType", "DELETE",
		"endpoint", "/api/main-page-entries/"+mainPageIDStr,
	)

	mainPageID, err := strconv.Atoi(mainPageIDStr)
	if err != nil {
		log.Warn("invalid removal path parameter index format syntax errors flagged", "raw_id", mainPageIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid main page entry ID"})
		return
	}

	if err := h.service.DeleteMainPageEntry(mainPageID); err != nil {
		log.Error("failed to drop isolated configuration slot mapping data configurations rows", "entry_id", mainPageID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("dashboard landing layout shortcut slots completely drop sequence actions confirmed", "deleted_entry_id", mainPageID)
	c.JSON(http.StatusOK, gin.H{"message": "Main page entry deleted successfully"})
}

// UpdateMainPageEntry updates position or appearance of a main page entry
// @Summary Update main page entry (Admin only)
// @Description Modifies position or meta layout details of an entry on the main page by its ID. Requires administrator role.
// @Tags main-page
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Main Page Entry ID"
// @Param entry body model.UpdateMainPageEntryRequest true "Main page entry values"
// @Success 200 {object} map[string]string "Main page entry updated successfully"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden - Admin only"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/main-page-entries/{id} [put]
func (h *ProjectHandler) UpdateMainPageEntry(c *gin.Context) {
	mainPageIDStr := c.Param("id")
	log := logger.Log.With(
		"requestType", "PUT",
		"endpoint", "/api/main-page-entries/"+mainPageIDStr,
	)

	mainPageID, err := strconv.Atoi(mainPageIDStr)
	if err != nil {
		log.Warn("invalid query criteria specification conversion argument pattern mismatch", "raw_id", mainPageIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid main page entry ID"})
		return
	}

	var req model.UpdateMainPageEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to translate dashboard slots adjustment properties fields values", "entry_id", mainPageID, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMainPageEntry(mainPageID, req.Position, req.IconURL); err != nil {
		log.Error("failed to override parameters configuration bounds variables schemas inside table", "entry_id", mainPageID, "target_position", req.Position, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("dashboard layout parameter attributes completely synced and updated successfully", "entry_id", mainPageID, "new_position", req.Position)
	c.JSON(http.StatusOK, gin.H{"message": "Main page entry updated successfully"})
}