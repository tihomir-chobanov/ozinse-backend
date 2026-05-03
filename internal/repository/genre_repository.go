package repository

import (
	"database/sql"
	"fmt"
	"ozinse-backend/internal/model"
)

type GenreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{db: db}
}

func (r *GenreRepository) GetAll() ([]model.Genre, error) {
	rows, err := r.db.Query(`SELECT id, name, icon_url FROM genre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []model.Genre
	for rows.Next() {
		var g model.Genre
		if err := rows.Scan(&g.ID, &g.Name, &g.IconUrl); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func (r *GenreRepository) GetByID(id int) (*model.Genre, error) {
	var g model.Genre
	err := r.db.QueryRow(`SELECT id, name, icon_url FROM genre WHERE id = $1`, id).
		Scan(&g.ID, &g.Name, &g.IconUrl)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GenreRepository) Create(g *model.Genre) error {
	return r.db.QueryRow(
		`INSERT INTO genre (name, icon_url) VALUES ($1, $2) RETURNING id`,
		g.Name, g.IconUrl,
	).Scan(&g.ID)
}

func (r *GenreRepository) Update(g *model.Genre) error {
	// We use Exec to update the genre
	result, err := r.db.Exec(
		`UPDATE genre SET name = $1, icon_url = $2 WHERE id = $3`,
		g.Name, g.IconUrl, g.ID,
	)
	if err != nil {
		return err
	}

	// We check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, we return an error
	if rowsAffected == 0 {
		return fmt.Errorf("genre with id %d not found", g.ID)
	}

	return nil
}

func (r *GenreRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM genre WHERE id = $1`, id)
	if err != nil {
		return err
	}
	// if we didn't implement RowsAffected(), we could delete a genre that doesn't exist and that is missleading/wrong
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("genre with id %d not found", id)
	}
	return nil
}

func (r *GenreRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM genre WHERE name = $1)`
	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

/* The Repository acts as a bridge between your application's business logic and the data source (database). It is a design pattern that isolates data access logic. Instead of your code knowing exactly how to write SQL queries, it simply asks the repository for an object (e.g., "Give me Category with ID 5"). The repository handles the technical execution with the database and returns a clean Model.

Repository (The Supplier): The only layer allowed to communicate with the database (PostgreSQL) using raw SQL queries.

*/
