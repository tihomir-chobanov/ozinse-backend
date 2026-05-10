package model

import "time"

// User represents a system user account.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // The password is hidden from JSON responses for security
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	BirthDate time.Time `json:"birth_date"`
	RoleID    int       `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	Image     string    `json:"image"`
}