package service

import (
	"fmt"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetByID(id int) (*model.Category, error) {
	return s.repo.GetByID(id)
}

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

func (s *CategoryService) Update(c *model.Category) error {
	return s.repo.Update(c)
}

func (s *CategoryService) Delete(id int) error {
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
