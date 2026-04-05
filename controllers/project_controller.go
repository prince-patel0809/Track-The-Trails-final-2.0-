package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"track-the-trails/config"
	"track-the-trails/middlewares"
	"track-the-trails/models"
)

// ================= CREATE PROJECT =================
func CreateProject(w http.ResponseWriter, r *http.Request) {

	// ✅ GET USER FROM CONTEXT (YOUR CODE)
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.CreateProjectInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Name == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var projectID int
	err = config.DB.QueryRow(
		"INSERT INTO projects (name, description, created_by) VALUES ($1,$2,$3) RETURNING project_id",
		input.Name, input.Description, userID,
	).Scan(&projectID)

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Project created",
		"project_id": projectID,
	})
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {

	// ===== GET USER FROM CONTEXT =====
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ===== GET PROJECT ID FROM URL =====
	idStr := strings.TrimPrefix(r.URL.Path, "/project/delete/")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// ===== CHECK OWNER =====
	var creatorID int
	err = config.DB.QueryRow(
		"SELECT created_by FROM projects WHERE project_id=$1",
		projectID,
	).Scan(&creatorID)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// ===== ONLY OWNER CAN DELETE =====
	if creatorID != userID {
		http.Error(w, "Forbidden: Only owner can delete", http.StatusForbidden)
		return
	}

	// ===== DELETE PROJECT =====
	_, err = config.DB.Exec(
		"DELETE FROM projects WHERE project_id=$1",
		projectID,
	)

	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	// ===== SUCCESS RESPONSE =====
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Project deleted successfully",
	})
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {

	// ===== GET USER FROM CONTEXT =====
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ===== GET PROJECT ID =====
	idStr := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// ===== PARSE INPUT =====
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Name == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// ===== CHECK OWNER =====
	var creatorID int
	err = config.DB.QueryRow(
		"SELECT created_by FROM projects WHERE project_id=$1",
		projectID,
	).Scan(&creatorID)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// ===== ONLY CREATOR CAN UPDATE =====
	if creatorID != userID {
		http.Error(w, "Forbidden: Only owner can update", http.StatusForbidden)
		return
	}

	// ===== UPDATE PROJECT =====
	_, err = config.DB.Exec(
		"UPDATE projects SET name=$1, description=$2 WHERE project_id=$3",
		input.Name, input.Description, projectID,
	)

	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	// ===== RESPONSE =====
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Project updated successfully",
	})
}
