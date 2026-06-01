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

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=7"`
}

// ForgotPasswordRequest defines the input when a user requests a reset link.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest defines the input when the user submits their new password.
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
