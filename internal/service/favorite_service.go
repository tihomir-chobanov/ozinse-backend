package service

import "ozinse-backend/internal/repository"

type FavoriteService struct {
	repo *repository.FavoriteRepository
}

func NewFavoriteService(repo *repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) AddFavorite(userID, projectID int) error {
	return s.repo.AddFavorite(userID, projectID)
}

func (s *FavoriteService) RemoveFavorite(userID, projectID int) error {
	return s.repo.RemoveFavorite(userID, projectID)
}