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

func GetMyProjects(w http.ResponseWriter, r *http.Request) {

	// ===== GET USER FROM CONTEXT =====
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ===== FETCH PROJECTS =====
	rows, err := config.DB.Query(
		`SELECT project_id, name, description, created_at 
		 FROM projects 
		 WHERE created_by = $1 
		 ORDER BY created_at DESC`,
		userID,
	)

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	// ===== STORE RESULTS =====
	var projects []map[string]interface{}

	for rows.Next() {
		var id int
		var name, description, createdAt string

		err := rows.Scan(&id, &name, &description, &createdAt)
		if err != nil {
			http.Error(w, "Error reading data", 500)
			return
		}

		projects = append(projects, map[string]interface{}{
			"project_id":  id,
			"name":        name,
			"description": description,
			"created_at":  createdAt,
		})
	}

	// ===== RESPONSE =====
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func GetProjectDetails(w http.ResponseWriter, r *http.Request) {

	// ===== GET PROJECT ID FROM URL =====
	parts := strings.Split(r.URL.Path, "/")
	projectID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// ===== QUERY DATABASE =====
	var project struct {
		ProjectID   int    `json:"project_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		OwnerID     int    `json:"owner_id"`
		OwnerName   string `json:"owner_name"`
		OwnerEmail  string `json:"owner_email"`
	}

	err = config.DB.QueryRow(`
		SELECT 
			p.project_id,
			p.name,
			p.description,
			p.created_at,
			u.user_id,
			u.name,
			u.email
		FROM projects p
		JOIN users u ON p.created_by = u.user_id
		WHERE p.project_id = $1
	`, projectID).Scan(
		&project.ProjectID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.OwnerID,
		&project.OwnerName,
		&project.OwnerEmail,
	)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// ===== RESPONSE =====
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func AddMember(w http.ResponseWriter, r *http.Request) {

	// ===== GET USER FROM CONTEXT =====
	val := r.Context().Value(middlewares.UserIDKey)
	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ===== GET PROJECT ID =====
	parts := strings.Split(r.URL.Path, "/")
	projectID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// ===== REQUEST BODY =====
	var input struct {
		UserID int    `json:"user_id"`
		Role   string `json:"role"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.UserID == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Role == "" {
		input.Role = "member"
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

	if creatorID != userID {
		http.Error(w, "Only owner can add members", http.StatusForbidden)
		return
	}

	// ===== INSERT MEMBER =====
	_, err = config.DB.Exec(
		"INSERT INTO project_members (project_id, user_id, role) VALUES ($1,$2,$3)",
		projectID, input.UserID, input.Role,
	)

	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Member added successfully",
	})
}
