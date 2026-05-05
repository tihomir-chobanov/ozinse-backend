package service

import (
	"fmt"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

// GenreService handles business logic for genre operations.
type GenreService struct {
	repo *repository.GenreRepository
}

// NewGenreService creates and returns a new GenreService instance.
func NewGenreService(repo *repository.GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

// GetAll retrieves all genres from the repository.
func (s *GenreService) GetAll() ([]model.Genre, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a genre by its ID.
func (s *GenreService) GetByID(id int) (*model.Genre, error) {
	return s.repo.GetByID(id)
}

// Create creates a new genre after checking if one with the same name already exists.
func (s *GenreService) Create(g *model.Genre) error {
	// 1. Check if genre with name already exists
	exists, err := s.repo.ExistsByName(g.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("genre with name '%s' already exists", g.Name)
	}
	// 2. If not, create genre
	return s.repo.Create(g)
}

// Update updates an existing genre.
func (s *GenreService) Update(c *model.Genre) error {
	return s.repo.Update(c)
}

// Delete removes a genre by its ID.
func (s *GenreService) Delete(id int) error {
	return s.repo.Delete(id)
}
