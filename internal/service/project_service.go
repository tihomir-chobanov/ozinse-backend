package service

import (
	"fmt"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

// ProjectService handles business logic for project operations.
type ProjectService struct {
	repo *repository.ProjectRepository
}

// NewProjectService creates and returns a new ProjectService instance.
func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// GetAll retrieves a paginated list of projects along with the total count of records.
// It calculates the database offset based on the requested page and limit parameters.
func (s *ProjectService) GetAll(page int, limit int, search string) ([]model.Project, int, error) {
	// Sanitize inputs to prevent negative or zero values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	// Calculate how many records PostgreSQL needs to skip
	// Example: Page 2 with Limit 10 -> Offset = (2 - 1) * 10 = 10 (skips first 10 records)
	offset := (page - 1) * limit

	// Call the repository layer with calculated pagination parameters
	return s.repo.GetAll(limit, offset, search)
}

// GetByID retrieves a project by its ID.
func (s *ProjectService) GetByID(id int) (*model.Project, error) {
	return s.repo.GetByID(id)
}

// Create creates a new project with associated genres, age categories, and categories after checking if one with the same title already exists.
func (s *ProjectService) Create(p *model.Project, genreIDs []int, ageCategoryIDs []int, categoryIDs []int) error {
	// 1. Check if project's title already exists
	exists, err := s.repo.ExistsByName(p.Title)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("project with name '%s' already exists", p.Title)
	}
	// 2. If not, create project
	return s.repo.Create(p, genreIDs, ageCategoryIDs, categoryIDs)
}

// Update updates an existing project.
func (s *ProjectService) Update(c *model.Project) error {
	return s.repo.Update(c)
}

// Delete removes a project by its ID.
func (s *ProjectService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProjectService) GetTrending() ([]model.Project, error) {
    return s.repo.GetTrendingProjects()
}

func (s *ProjectService) GetSimilar(projectID int) ([]model.Project, error) {
	return s.repo.GetSimilarProjects(projectID)
}