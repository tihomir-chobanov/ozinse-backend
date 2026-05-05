package repository

import (
	"database/sql"
	"fmt"
	"ozinse-backend/internal/model"
)

// CategoryRepository provides CRUD access to the category table.
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new CategoryRepository instance.
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll retrieves all categories from the database.
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

// GetByID retrieves a category by its unique ID.
func (r *CategoryRepository) GetByID(id int) (*model.Category, error) {
	var c model.Category
	err := r.db.QueryRow(`SELECT id, name FROM category WHERE id = $1`, id).
		Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Create inserts a new category and returns the generated ID.
func (r *CategoryRepository) Create(c *model.Category) error {
	return r.db.QueryRow(
		`INSERT INTO category (name) VALUES ($1) RETURNING id`,
		c.Name,
	).Scan(&c.ID)
}

// Update modifies an existing category by ID. Returns an error if no rows were affected.
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

// Delete removes a category from the database by ID.
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

// ExistsByName checks whether a category with the provided name already exists.
func (r *CategoryRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM category WHERE name = $1)`
	err := r.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
