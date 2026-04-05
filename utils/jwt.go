package utils

import (
	"encoding/base64"
	"fmt"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getKey() []byte {
	secret := os.Getenv("JWT_SECRET")

	key, err := base64.StdEncoding.DecodeString(secret)

	if err != nil {
		panic("Invalid JWT secret (base64 decode failed)")
	}

	return key
}

// generate token

func GenerateToken(userID int, email string) (string, error) {

	claim := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(getKey())
}

// verify token

func VerifyToken(tokenStr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return getKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
