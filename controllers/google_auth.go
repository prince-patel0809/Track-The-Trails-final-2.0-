package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"track-the-trails/config"
	"track-the-trails/utils"
)

type GoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {

	url := utils.GoogleOAuthConfig.AuthCodeURL("state-token")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {

	// Get auth code from Google
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := utils.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	client := utils.GoogleOAuthConfig.Client(context.Background(), token)

	// Get Google user info
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser GoogleUser
	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		http.Error(w, "Failed to decode Google user", http.StatusInternalServerError)
		return
	}

	// Check if user exists
	var userID int
	err = config.DB.QueryRow(
		"SELECT user_id FROM users WHERE email=$1",
		googleUser.Email,
	).Scan(&userID)

	// If user does not exist, create new user
	if err == sql.ErrNoRows {

		err = config.DB.QueryRow(`
			INSERT INTO users(name,email,password,role)
			VALUES($1,$2,$3,$4)
			RETURNING user_id
		`,
			googleUser.Name,
			googleUser.Email,
			"",
			"member",
		).Scan(&userID)

		if err != nil {
			http.Error(w, "User creation failed", http.StatusInternalServerError)
			return
		}

	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Generate JWT
	jwtToken, err := utils.GenerateToken(userID, googleUser.Email)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// Redirect back to Android app
	redirectURL := "myapp://auth-success?token=" + jwtToken

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
