package service

import (
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for user operations.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService initializes a new UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}


func (s *UserService) GetAll() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetByID(id int) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(user *model.User) error {
	// If the password is being updated, we need to hash it before saving
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return s.repo.Update(user)
}

func (s *UserService) Delete(id int) error {
	return s.repo.Delete(id)
}