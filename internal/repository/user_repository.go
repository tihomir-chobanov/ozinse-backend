package repository

import (
	"database/sql"
	"errors"
	"ozinse-backend/internal/model"
)

// UserRepository handles database operations related to users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail retrieves a user from the database by their email address.
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	
	// Define the raw SQL query. 
query := `
		SELECT 
			id, email, password, full_name, phone, birth_date, role_id, created_at, image 
		FROM users 
		WHERE email = $1
	`
	
	// Execute the query and scan the result into the corresponding struct fields.
	// The order of arguments in Scan MUST match the order in the SELECT statement.
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.Phone,
		&user.BirthDate,
		&user.RoleID,
		&user.CreatedAt,
		&user.Image,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}