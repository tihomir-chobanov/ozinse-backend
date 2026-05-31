package repository

import (
	"database/sql"
	"fmt"
	
)

// FavoriteRepository provides CRUD access to the user_favorites table.
type FavoriteRepository struct {
	db *sql.DB
}


// NewFavoriteRepository creates a new instance of FavoriteRepository.
func NewFavoriteRepository(db *sql.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) AddFavorite(userID, projectID int) error {
	query := `
		INSERT INTO user_favorites (user_id, project_id) 
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`
	
	res, err := r.db.Exec(query, userID, projectID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return nil 
	}

	return nil
}

func (r *FavoriteRepository) RemoveFavorite(userID, projectID int) error {
	query := `DELETE FROM user_favorites WHERE user_id = $1 AND project_id = $2`
	res, err := r.db.Exec(query, userID, projectID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("project %d is not in user %d favorites", projectID, userID)
	}

	return nil
}