package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"track-the-trails/config"
	"track-the-trails/models"
	"track-the-trails/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Missing", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	var userID int
	err = config.DB.QueryRow(
		`INSERT INTO Users
		(name,email,password,role)
		VALUES ($1,$2,$3,$4)
		RETURNING user_id`,
		user.Name,
		user.Email,
		string(hashed),

		"member",
	).Scan(&userID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			http.Error(w, "Email already exists", 400)
		} else {
			http.Error(w, "Database error", 500)
		}
		return
	}

	token, err := utils.GenerateToken(userID, user.Email)

	user.Password = ""

	response := map[string]interface{}{
		"message": "Registration successful",
		"token":   token,
		"user":    user,
	}

	json.NewEncoder(w).Encode(response)

}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	var login models.LoginInput

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if login.Email == "" || login.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	var userID int
	var name string
	var email string
	var hashedPassword string

	err = config.DB.QueryRow(
		`SELECT user_id, name, email, password
		FROM Users
		WHERE LOWER(email)=LOWER($1)`,
		login.Email,
	).Scan(&userID, &name, &email, &hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(login.Password),
	)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(userID, email)
	if err != nil {
		http.Error(w, "Token error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user": map[string]interface{}{
			"user_id": userID,
			"name":    name,
			"email":   email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	rows, err := config.DB.Query(`
		SELECT user_id, name, email, profile_image, role, created_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		UserID       int       `json:"user_id"`
		Name         string    `json:"name"`
		Email        string    `json:"email"`
		ProfileImage string    `json:"profile_image"`
		Role         string    `json:"role"`
		CreatedAt    time.Time `json:"created_at"`
	}

	var users []User

	for rows.Next() {

		var user User
		var profileImage sql.NullString // ✅ FIX

		err := rows.Scan(
			&user.UserID,
			&user.Name,
			&user.Email,
			&profileImage, // ✅ use temp var
			&user.Role,
			&user.CreatedAt,
		)

		if err != nil {
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}

		// ✅ HANDLE NULL VALUE
		if profileImage.Valid {
			user.ProfileImage = profileImage.String
		} else {
			user.ProfileImage = ""
		}

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
