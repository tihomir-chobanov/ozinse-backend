package service

import (
	"fmt"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

// AgeCategoryService handles business logic for category operations.
type AgeCategoryService struct {
	repo *repository.AgeCategoryRepository
}

// NewAgeCategoryService creates and returns a new AgeCategoryService instance.
func NewAgeCategoryService(repo *repository.AgeCategoryRepository) *AgeCategoryService {
	return &AgeCategoryService{repo: repo}
}

// GetAll retrieves all Agecategories from the repository.
func (s *AgeCategoryService) GetAll() ([]model.Age_Category, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a Agecategory by its ID.
func (s *AgeCategoryService) GetByID(id int) (*model.Age_Category, error) {
	return s.repo.GetByID(id)
}

// Create creates a new agecategory after checking if one with the same name already exists.
func (s *AgeCategoryService) Create(ac *model.Age_Category) error {
	exists, err := s.repo.ExistsByName(ac.Range)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Agecategory with name '%s' already exists", ac.Range)
	}
	return s.repo.Create(ac)
}

// Update updates an existing category.
func (s *AgeCategoryService) Update(c *model.Age_Category) error {
	return s.repo.Update(c)
}

// Delete removes a category by its ID.
func (s *AgeCategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
