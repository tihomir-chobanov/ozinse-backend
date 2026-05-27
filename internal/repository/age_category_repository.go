package repository

import (
	"database/sql"
	"fmt"
	"ozinse-backend/internal/model"
)

// AgeCategoryRepository provides CRUD access to the age_category table.
type AgeCategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new CategoryRepository instance.
func NewAgeCategoryRepository(db *sql.DB) *AgeCategoryRepository {
	return &AgeCategoryRepository{db: db}
}

// GetAll retrieves all age_categories from the database.
func (r *AgeCategoryRepository) GetAll() ([]model.Age_Category, error) {
	rows, err := r.db.Query(`SELECT id, range, icon_url FROM age_category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var age_categories []model.Age_Category
	for rows.Next() {
		var ac model.Age_Category
		if err := rows.Scan(&ac.ID, &ac.Range, &ac.IconUrl); err != nil {
			return nil, err
		}
		age_categories = append(age_categories, ac)
	}
	return age_categories, nil
}

// GetByID retrieves a age_category by its unique ID.
func (r *AgeCategoryRepository) GetByID(id int) (*model.Age_Category, error) {
	var ac model.Age_Category
	err := r.db.QueryRow(`SELECT id, range, icon_url FROM age_category WHERE id = $1`, id).
		Scan(&ac.ID, &ac.Range, &ac.IconUrl)
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

// Create inserts a new age_category and returns the generated ID.
func (r *AgeCategoryRepository) Create(c *model.Age_Category) error {
	return r.db.QueryRow(
		`INSERT INTO age_category (range, icon_url) VALUES ($1, $2) RETURNING id`,
		c.Range, c.IconUrl,
	).Scan(&c.ID)
}

// Update modifies an existing age_category by ID. Returns an error if no rows were affected.
func (r *AgeCategoryRepository) Update(c *model.Age_Category) error {
	// We use Exec to update the category
	result, err := r.db.Exec(
		`UPDATE age_category SET range = $1, icon_url = $2 WHERE id = $3`,
		c.Range, c.IconUrl, c.ID,
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
		return fmt.Errorf("age_category with id %d not found", c.ID)
	}

	return nil
}

// Delete removes a age_category from the database by ID.
func (r *AgeCategoryRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM age_category WHERE id = $1`, id)
	if err != nil {
		return err
	}
	// if we didn't implement RowsAffected(), we could delete a age_category that doesn't exist and that is missleading/wrong
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("age_category with id %d not found", id)
	}
	return nil
}

// ExistsByName checks whether a category with the provided name already exists.
func (r *AgeCategoryRepository) ExistsByName(rangeVal string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM age_category WHERE range = $1)`
	err := r.db.QueryRow(query, rangeVal).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
