package repository

import (
	"database/sql"
	"fmt"
	"ozinse-backend/internal/model"
	"time"
)

// UserRepository handles database operations related to users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail retrieves a user record from the database by their unique email address.
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, email, password, full_name, phone, birth_date, role_id, created_at, image, language, notifications_enabled, dark_mode_enabled
		FROM users 
		WHERE email = $1`

	rows, err := r.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// If no row is found, Next() will be false, and we safely return nil, nil
	if !rows.Next() {
		return nil, nil
	}

	var user model.User
	err = rows.Scan(
		&user.ID, &user.Email, &user.Password, &user.FullName,
		&user.Phone, &user.BirthDate, &user.RoleID, &user.CreatedAt, &user.Image, &user.Language, &user.NotificationsEnabled, &user.DarkModeEnabled,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create inserts a new user record into the users table.
func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (email, password, full_name, phone, birth_date, role_id, image, language, notifications_enabled, dark_mode_enabled) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id, created_at`

	// Execute the insert statement and scan the database-generated ID and CreatedAt fields
	err := r.db.QueryRow(
		query,
		user.Email, user.Password, user.FullName,
		user.Phone, user.BirthDate, user.RoleID, user.Image, user.Language, user.NotificationsEnabled, user.DarkModeEnabled,
	).Scan(&user.ID, &user.CreatedAt)

	return err
}

// SaveResetToken stores the password reset token and its expiration time for a user.
func (r *UserRepository) SaveResetToken(email string, token string, expiresAt time.Time) error {
	query := `UPDATE users SET reset_token = $1, reset_token_expires_at = $2 WHERE email = $3`
	_, err := r.db.Exec(query, token, expiresAt, email)
	return err
}

// GetByResetToken finds a user by a valid reset token.
func (r *UserRepository) GetByResetToken(token string) (*model.User, error) {
	query := `
		SELECT id, email, password, full_name, phone, birth_date, role_id, reset_token, reset_token_expires_at 
		FROM users 
		WHERE reset_token = $1`

	var user model.User
	var resetToken sql.NullString
	var expiresAt sql.NullTime

	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.Email, &user.Password, &user.FullName,
		&user.Phone, &user.BirthDate, &user.RoleID, &resetToken, &expiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Check if the token has expired compared to the current time
	if expiresAt.Valid && expiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return &user, nil
}

// UpdatePassword updates only the password field and clears the reset token.
func (r *UserRepository) UpdatePassword(id int, hashedPassword string) error {
	query := `UPDATE users SET password = $1, reset_token = NULL, reset_token_expires_at = NULL WHERE id = $2`
	_, err := r.db.Exec(query, hashedPassword, id)
	return err
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	query := `SELECT id, email, full_name, phone, birth_date, role_id, created_at, image, language, notifications_enabled, dark_mode_enabled FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.ID, &user.Email, &user.FullName, &user.Phone, &user.BirthDate, &user.RoleID, &user.CreatedAt, &user.Image, &user.Language, &user.NotificationsEnabled, &user.DarkModeEnabled,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	query := `
		SELECT id, email, password, full_name, phone, birth_date, role_id, created_at, image, language, notifications_enabled, dark_mode_enabled 
		FROM users 
		WHERE id = $1`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	var user model.User
	err = rows.Scan(
		&user.ID, &user.Email, &user.Password, &user.FullName,
		&user.Phone, &user.BirthDate, &user.RoleID, &user.CreatedAt, &user.Image, &user.Language, &user.NotificationsEnabled, &user.DarkModeEnabled,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	query := `
		UPDATE users 
		SET email = $1, full_name = $2, phone = $3, birth_date = $4, role_id = $5, image = $6, language = $7, notifications_enabled = $8, dark_mode_enabled = $9
		WHERE id = $10`
	
	res, err := r.db.Exec(query, user.Email, user.FullName, user.Phone, user.BirthDate, user.RoleID, user.Image, user.Language, user.NotificationsEnabled, user.DarkModeEnabled, user.ID)
	if err != nil {
		return err 
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d to update", user.ID)
	}

	return nil
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err 
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d to delete", id)
	}

	return nil
}
