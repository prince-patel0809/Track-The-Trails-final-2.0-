package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Get user id from JWT token
func GetUserIDFromToken(r *http.Request) (int, error) {

	// ===== GET AUTH HEADER =====
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("authorization header missing")
	}

	// Expected format: Bearer TOKEN
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return 0, errors.New("invalid token format")
	}

	tokenString := parts[1]

	// ===== PARSE TOKEN =====
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getKey(), nil // your JWT secret key function
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	// ===== EXTRACT CLAIMS =====
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	// ===== GET USER ID =====
	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("user id not found")
	}

	userID := int(idFloat)

	return userID, nil
}
