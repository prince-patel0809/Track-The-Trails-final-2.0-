package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"track-the-trails/config"
	"track-the-trails/middlewares"
	"track-the-trails/models"
)

func extractID(r *http.Request) (int, error) {
	// Remove leading & trailing "/" then split
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Get last part (this should be the ID)
	idStr := strings.TrimSpace(parts[len(parts)-1])

	// Convert string → int
	return strconv.Atoi(idStr)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {

	var todo models.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	// ✅ FIXED TYPE (int)
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		fmt.Println("❌ VERIFY ERROR:", err)
		http.Error(w, "Unauthorized", 401)
		return
	}

	err = config.DB.QueryRow(`
	INSERT INTO todo(
		user_id,
		title,
		description,
		priority,
		due_date,
		due_time
	)
	VALUES($1,$2,$3,$4,$5,$6)
	RETURNING todo_id
	`,
		userID,
		todo.Title,
		todo.Description,
		todo.Priority,
		todo.DueDate,
		todo.DueTime,
	).Scan(&todo.TodoID)

	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "DB error", 500)
		return
	}

	todo.UserID = userID
	todo.Status = "pending"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo

	// ✅ GET USER FROM CONTEXT (NOT TOKEN AGAIN)
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := config.DB.Query(`
		SELECT 
			todo_id,
			user_id,
			title,
			description,
			status,
			priority,
			due_date,
			due_time
		FROM todo
		WHERE user_id = $1
	`, userID)

	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo

		err := rows.Scan(
			&todo.TodoID,
			&todo.UserID,
			&todo.Title,
			&todo.Description,
			&todo.Status,
			&todo.Priority,
			&todo.DueDate,
			&todo.DueTime,
		)

		if err != nil {
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}

		todos = append(todos, todo)
	}

	// ✅ IMPORTANT
	if err = rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	// ✅ GET ID FROM URL
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimSpace(parts[len(parts)-1])
	fmt.Println("ID STRING:", idStr)

	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Todo ID: "+idStr, http.StatusBadRequest)
		return
	}

	// ✅ DECODE BODY
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("DECODED TODO:", todo)

	// ✅ GET USER FROM CONTEXT (IMPORTANT FIX)
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ✅ UPDATE QUERY
	result, err := config.DB.Exec(`
		UPDATE todo
		SET title=$1, description=$2, status=$3, priority=$4, due_date=$5, due_time=$6
		WHERE todo_id=$7 AND user_id=$8
	`,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.Priority,
		todo.DueDate,
		todo.DueTime,
		todoID,
		userID,
	)

	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Todo not found or unauthorized", http.StatusNotFound)
		return
	}

	// ✅ RESPONSE
	todo.TodoID = todoID
	todo.UserID = userID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {

	// ✅ EXTRACT ID FROM URL
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimSpace(parts[len(parts)-1])
	fmt.Println("ID STRING:", idStr)

	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
		return
	}

	// ✅ GET USER FROM CONTEXT (FIXED)
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ✅ DELETE QUERY (SECURE)
	result, err := config.DB.Exec(`
		DELETE FROM todo
		WHERE todo_id = $1 AND user_id = $2
	`, todoID, userID)

	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ CHECK IF DELETED
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking delete", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Todo not found or unauthorized", http.StatusNotFound)
		return
	}

	// ✅ SUCCESS RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo deleted successfully",
	})
}

func CompleteTodo(w http.ResponseWriter, r *http.Request) {

	// ✅ EXTRACT ID FROM URL
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimSpace(parts[len(parts)-1])
	fmt.Println("ID STRING:", idStr)

	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
		return
	}

	// ✅ GET USER FROM CONTEXT (FIXED)
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ✅ UPDATE STATUS ONLY
	result, err := config.DB.Exec(`
		UPDATE todo
		SET status = 'completed'
		WHERE todo_id = $1 AND user_id = $2
	`, todoID, userID)

	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Todo not found or unauthorized", http.StatusNotFound)
		return
	}

	// ✅ SUCCESS RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo marked as completed",
	})
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {

	// ✅ GET ID FROM URL
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimSpace(parts[len(parts)-1])
	fmt.Println("ID STRING:", idStr)

	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
		return
	}

	// ✅ GET USER FROM CONTEXT (IMPORTANT)
	val := r.Context().Value(middlewares.UserIDKey)

	userID, ok := val.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ✅ QUERY ONE TODO
	var todo models.Todo

	err = config.DB.QueryRow(`
		SELECT 
			todo_id,
			user_id,
			title,
			description,
			status,
			priority,
			due_date,
			due_time
		FROM todo
		WHERE todo_id = $1 AND user_id = $2
	`, todoID, userID).Scan(
		&todo.TodoID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.Priority,
		&todo.DueDate,
		&todo.DueTime,
	)

	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// ✅ RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}
