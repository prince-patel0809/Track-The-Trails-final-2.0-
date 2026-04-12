package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"track-the-trails/config"
	"track-the-trails/middlewares"
	"track-the-trails/utils"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {

	// ===== GET USER FROM TOKEN =====
	val := r.Context().Value(middlewares.UserIDKey)
	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ===== GET PROJECT ID =====
	parts := strings.Split(r.URL.Path, "/")
	projectID, err := strconv.Atoi(parts[len(parts)-2])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// ===== CHECK OWNER =====
	var ownerID int
	err = config.DB.QueryRow(
		"SELECT created_by FROM projects WHERE project_id=$1",
		projectID,
	).Scan(&ownerID)

	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Only project owner can assign task", http.StatusForbidden)
		return
	}

	// ===== REQUEST BODY =====
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AssignedTo  int    `json:"assigned_to"`
		Priority    string `json:"priority"`
		DueDate     string `json:"due_date"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Title == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// ===== CHECK ASSIGNED USER IS MEMBER =====
	var exists int
	err = config.DB.QueryRow(
		`SELECT 1 FROM project_members 
		 WHERE project_id=$1 AND user_id=$2`,
		projectID, input.AssignedTo,
	).Scan(&exists)

	if err != nil {
		http.Error(w, "User is not a project member", http.StatusBadRequest)
		return
	}

	// ===== INSERT TASK =====
	var taskID int
	err = config.DB.QueryRow(`
		INSERT INTO tasks 
		(project_id, title, description, assigned_to, created_by, priority)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING task_id
	`,
		projectID,
		input.Title,
		input.Description,
		input.AssignedTo,
		userID,
		input.Priority,
	).Scan(&taskID)

	if err != nil {
		http.Error(w, "Failed to create task", 500)
		return
	}

	// ===== GET EMAIL =====
	var email string
	err = config.DB.QueryRow(
		"SELECT email FROM users WHERE user_id=$1",
		input.AssignedTo,
	).Scan(&email)

	// ===== GET PROJECT NAME =====
	var projectName string
	config.DB.QueryRow(
		"SELECT name FROM projects WHERE project_id=$1",
		projectID,
	).Scan(&projectName)

	// ===== GET OWNER NAME =====
	var ownerName string
	config.DB.QueryRow(
		"SELECT name FROM users WHERE user_id=$1",
		userID,
	).Scan(&ownerName)

	// ===== SEND EMAIL =====
	if err == nil {
		go utils.SendTaskAssignmentEmail(
			email,
			projectName,
			ownerName,
			input.Title,
			input.Description,
		)
	}
	log.Println("Sending email to:", email)

	// ===== RESPONSE =====
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Task assigned successfully",
		"task_id": taskID,
	})
}

func GetTasksByProject(w http.ResponseWriter, r *http.Request) {

	// ===== GET USER FROM TOKEN =====
	val := r.Context().Value(middlewares.UserIDKey)
	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ===== GET PROJECT ID FROM URL =====
	parts := strings.Split(r.URL.Path, "/")
	projectID, err := strconv.Atoi(parts[len(parts)-2])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// ===== CHECK ACCESS (OWNER OR MEMBER) =====
	var exists int
	err = config.DB.QueryRow(`
		SELECT 1 FROM project_members 
		WHERE project_id=$1 AND user_id=$2
		UNION
		SELECT 1 FROM projects 
		WHERE project_id=$1 AND created_by=$2
	`, projectID, userID).Scan(&exists)

	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// ===== QUERY TASKS =====
	rows, err := config.DB.Query(`
		SELECT 
			t.task_id,
			t.title,
			t.description,
			t.status,
			t.priority,
			t.due_date,
			t.created_at,
			u.name
		FROM tasks t
		LEFT JOIN users u ON t.assigned_to = u.user_id
		WHERE t.project_id = $1
		ORDER BY t.created_at DESC
	`, projectID)

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	// ===== RESPONSE STRUCT =====
	type Task struct {
		TaskID       int     `json:"task_id"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		Status       string  `json:"status"`
		Priority     string  `json:"priority"`
		DueDate      *string `json:"due_date"`
		CreatedAt    string  `json:"created_at"`
		AssignedName string  `json:"assigned_name"`
	}

	var tasks []Task

	for rows.Next() {
		var t Task

		err := rows.Scan(
			&t.TaskID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.DueDate,
			&t.CreatedAt,
			&t.AssignedName,
		)

		if err != nil {
			http.Error(w, "Error reading data", 500)
			return
		}

		tasks = append(tasks, t)
	}

	// ===== RESPONSE =====
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"project_id": projectID,
		"tasks":      tasks,
	})
}
