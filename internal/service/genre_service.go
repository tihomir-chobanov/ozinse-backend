package service

import (
	"fmt"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

type GenreService struct {
	repo *repository.GenreRepository
}

func NewGenreService(repo *repository.GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) GetAll() ([]model.Genre, error) {
	return s.repo.GetAll()
}

func (s *GenreService) GetByID(id int) (*model.Genre, error) {
	return s.repo.GetByID(id)
}

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

func (s *GenreService) Update(c *model.Genre) error {
	return s.repo.Update(c)
}

func (s *GenreService) Delete(id int) error {
	return s.repo.Delete(id)
}

/*
In your project, the Service layer acts as the "Brain" or the "Chief Chef" of the application.
While the Handler manages the conversation with the user and
the Repository manages the database,
the Service focuses on the Business Logic.

What does the Service layer do?
Logic and Rules: For example, if a user wants to create a category, the Service could check if that category name already exists before telling the Repository to save it.

Orchestration: Sometimes a single request requires multiple database actions. The Service coordinates these. For instance, if you delete a category, the Service might first check if there are any genres or projects linked to it to prevent data errors.

Decoupling: It separates the "How" (Database/SQL) from the "Where" (HTTP/API). The Service doesn't care if you are using Gin or a command-line tool; it just cares about the rules of your app.

Summary of the Flow

Handler: Receives the request (e.g., "Delete ID 5").
Service: Receives the command from the Handler. It thinks: "Is it okay to delete this? Does ID 5 exist?".
Repository: Receives the instruction from the Service. It executes the SQL: DELETE FROM category WHERE id = 5.


*/
