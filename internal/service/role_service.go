package service

import (
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CreateRole(role *model.Role) error {
	return s.repo.Create(role)
}

func (s *RoleService) GetAllRoles() ([]model.Role, error) {
	return s.repo.GetAll()
}

func (s *RoleService) GetRoleByID(id int) (*model.Role, error) {
	return s.repo.GetByID(id)
}

func (s *RoleService) UpdateRole(role *model.Role) error {
	return s.repo.Update(role)
}	

func (s *RoleService) DeleteRole(id int) error {
	return s.repo.Delete(id)
}