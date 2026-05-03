package repository

import (
	"database/sql"
	"fmt"
	"ozinse-backend/internal/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	rows, err := r.db.Query(`SELECT id, name FROM category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*model.Category, error) {
	var c model.Category
	err := r.db.QueryRow(`SELECT id, name FROM category WHERE id = $1`, id).
		Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(c *model.Category) error {
	return r.db.QueryRow(
		`INSERT INTO category (name) VALUES ($1) RETURNING id`,
		c.Name,
	).Scan(&c.ID)
}

func (r *CategoryRepository) Update(c *model.Category) error {
	// We use Exec to update the category
	result, err := r.db.Exec(
		`UPDATE category SET name = $1 WHERE id = $2`,
		c.Name, c.ID,
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
		return fmt.Errorf("category with id %d not found", c.ID)
	}

	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM category WHERE id = $1`, id)
	if err != nil {
		return err
	}
	// if we didn't implement RowsAffected(), we could delete a category that doesn't exist and that is missleading/wrong
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with id %d not found", id)
	}
	return nil
}

func (r *CategoryRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM category WHERE name = $1)`
	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

/* The Repository acts as a bridge between your application's business logic and the data source (database). It is a design pattern that isolates data access logic. Instead of your code knowing exactly how to write SQL queries, it simply asks the repository for an object (e.g., "Give me Category with ID 5"). The repository handles the technical execution with the database and returns a clean Model.

Repository (The Supplier): The only layer allowed to communicate with the database (PostgreSQL) using raw SQL queries.

*/
