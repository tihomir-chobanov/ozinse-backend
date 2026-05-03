package service

import (
	"fmt"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetAll() ([]model.Project, error) {
	return s.repo.GetAll()
}

func (s *ProjectService) GetByID(id int) (*model.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectService) Create(p *model.Project) error {
	// 1. Check if project's title already exists
	exists, err := s.repo.ExistsByName(p.Title)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("project with name '%s' already exists", p.Title)
	}
	// 2. If not, create project
	return s.repo.Create(p)
}

func (s *ProjectService) Update(c *model.Project) error {
	return s.repo.Update(c)
}

func (s *ProjectService) Delete(id int) error {
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
