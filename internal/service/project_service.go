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

// GetAll retrieves all projects from the repository.
func (s *ProjectService) GetAll() ([]model.Project, error) {
	return s.repo.GetAll()
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
