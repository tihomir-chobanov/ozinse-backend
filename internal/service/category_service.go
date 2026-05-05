package service

import (
	"fmt"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

// CategoryService handles business logic for category operations.
type CategoryService struct {
	repo *repository.CategoryRepository
}

// NewCategoryService creates and returns a new CategoryService instance.
func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll retrieves all categories from the repository.
func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a category by its ID.
func (s *CategoryService) GetByID(id int) (*model.Category, error) {
	return s.repo.GetByID(id)
}

// Create creates a new category after checking if one with the same name already exists.
func (s *CategoryService) Create(c *model.Category) error {
	exists, err := s.repo.ExistsByName(c.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("category with name '%s' already exists", c.Name)
	}
	return s.repo.Create(c)
}

// Update updates an existing category.
func (s *CategoryService) Update(c *model.Category) error {
	return s.repo.Update(c)
}

// Delete removes a category by its ID.
func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
