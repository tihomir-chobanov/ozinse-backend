package service

import (
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
)

// UserService handles business logic for user operations.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService initializes a new UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetByEmail fetches a user by email using the repository layer.
func (s *UserService) GetByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}