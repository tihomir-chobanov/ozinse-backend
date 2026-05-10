package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the data that will be "sealed" inside the token.
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	RoleID int    `json:"role_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new encrypted JWT token for a successfully logged-in user.
func GenerateToken(userID int, email string, roleID int, secret string, expiryHours int) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Set the token expiration time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			// Set the token issuance time
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token itself using the HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign it with our secret key and return it as a string
	return token.SignedString([]byte(secret))
}

// ValidateToken checks if a given token is valid, not expired, and not tampered with.
func ValidateToken(tokenString string, secret string) (*Claims, error) {
	// Attempt to parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// If valid, return the claims data from it
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}