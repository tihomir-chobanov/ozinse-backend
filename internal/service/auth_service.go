package service

import (
	"errors"
	"fmt"
	"ozinse-backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"ozinse-backend/internal/repository"
	"time"
	"crypto/rand"
	"encoding/hex"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// GetByEmail fetches a user by email using the repository layer.
func (s *AuthService) GetByEmail(email string) (*model.User, error) {
	return s.userRepo.GetByEmail(email)
}

// Register handles the business logic for creating a new user profile with defaults.
func (s *AuthService) Register(req *model.RegisterRequest) (*model.User, error) {
	// 1. Check if a user with this email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("a user with this email address already exists")
	}

	// 2. Encrypt the raw password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// NEW: 3. Parse the default string date into a time.Time object
	// Go uses a unique layout standard based on the specific reference time: Mon Jan 2 15:04:05 MST 2006
	defaultBirthDate, err := time.Parse("2006-01-02", "1970-01-01")
	if err != nil {
		return nil, fmt.Errorf("failed to parse default birth date: %w", err)
	}

	// 4. Enrich the data with application defaults before saving
	newUser := &model.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  "New User",
		Phone:     "",
		BirthDate: defaultBirthDate, // <-- Now using the parsed time.Time object
		RoleID:    1,
		Image:     "user.png",
		Language:  "kk",
		NotificationsEnabled: true,
		DarkModeEnabled:      false,
	}

	// 5. Call repository layer to persist the data in PostgreSQL
	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	// We are hiding the password field in the response for security reasons, even though it is hashed
	newUser.Password = ""

	return newUser, nil
}

// GenerateRandomToken creates a secure random hex string to use as a reset token.
// This method is used in the ForgotPassword flow to create a unique token for password reset requests.
func generateRandomToken() string {
	b := make([]byte, 20)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ForgotPassword handles generating a secure token and setting its expiration (60 mins).
func (s *AuthService) ForgotPassword(req *model.ForgotPasswordRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		// Security tip: To prevent user enumeration, you can return a generic message 
		// instead of "user not found", but for development, an error helps.
		return "", errors.New("user with this email does not exist")
	}

	token := generateRandomToken()
	expiresAt := time.Now().Add(60 * time.Minute) // Token valid for 60 minutes

	err = s.userRepo.SaveResetToken(req.Email, token, expiresAt)
	if err != nil {
		return "", err
	}

	// In a real production app: s.emailService.SendResetEmail(req.Email, token)
	return token, nil
}

// ResetPassword validates the token and applies the new hashed password.
func (s *AuthService) ResetPassword(req *model.ResetPasswordRequest) error {
	user, err := s.userRepo.GetByResetToken(req.Token)
	if err != nil {
		return err // Will capture "token has expired" or DB errors
	}
	if user == nil {
		return errors.New("invalid password reset token")
	}

	// Hash the new password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password and clear the used token inside the database
	return s.userRepo.UpdatePassword(user.ID, string(hashedPassword))
}